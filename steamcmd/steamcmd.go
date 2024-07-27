package steamcmd

import (
	"bufio"
	"cs-server-controller/event"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/creack/pty"
)

var instanceCreated = false

type SteamcmdInstance struct {
	steamCmdDir string
	serverDir   string

	running  atomic.Bool
	canceled atomic.Bool
	lastLine atomic.Value

	pty *os.File
	cmd *exec.Cmd

	startStopLock sync.Mutex
	cmdLock       sync.Mutex

	onOutput    event.EventWithData[string]
	onStarted   event.Event
	onFinished  event.Event
	onCancelled event.Event
	onFailed    event.EventWithData[error]
}

func NewInstance(steamcmdDir, serverDir string, enableEventLogging bool) (*SteamcmdInstance, error) {
	if instanceCreated {
		return nil, errors.New("another instance already exists. Use only one instance throughout the program")
	}

	isValidSteamcmdDir, _ := govalidator.IsFilePath(steamcmdDir)
	if !isValidSteamcmdDir {
		errorMsg := fmt.Sprintf("steamcmd dir %q is not a valid filepath", steamcmdDir)
		return nil, errors.New(errorMsg)
	}

	isValidServerDir, _ := govalidator.IsFilePath(serverDir)
	if !isValidServerDir {
		errorMsg := fmt.Sprintf("server dir %q is not a valid filepath", serverDir)
		return nil, errors.New(errorMsg)
	}

	i := SteamcmdInstance{
		steamCmdDir: steamcmdDir,
		serverDir:   serverDir,
	}

	i.onFailed.Register(func(pwd event.PayloadWithData[error]) {
		i.cleanup()
		i.running.Store(false)
		i.lastLine.Store("")
	})

	i.onCancelled.Register(func(dp event.DefaultPayload) {
		i.cleanup()
		i.lastLine.Store("")
		i.running.Store(false)
		i.canceled.Store(false)
	})

	i.onFinished.Register(func(dp event.DefaultPayload) {
		i.cleanup()
		i.lastLine.Store("")
		i.running.Store(false)
	})

	if enableEventLogging {
		i.enableEventLogging()
	}

	instanceCreated = true
	return &i, nil
}

func (s *SteamcmdInstance) enableEventLogging() {
	s.onOutput.Register(func(pwd event.PayloadWithData[string]) {
		fmt.Println(pwd.TriggeredAtUtc.String() + " | steamcmd: " + pwd.Data)
	})

	s.onStarted.Register(func(dp event.DefaultPayload) {
		slog.Debug("onStarted", "triggeredAtUtc", dp.TriggeredAtUtc)
	})

	s.onFinished.Register(func(dp event.DefaultPayload) {
		slog.Debug("onFinished", "triggeredAtUtc", dp.TriggeredAtUtc)
	})

	s.onCancelled.Register(func(dp event.DefaultPayload) {
		slog.Debug("onStarted", "triggeredAtUtc", dp.TriggeredAtUtc)
	})

	s.onFailed.Register(func(pwd event.PayloadWithData[error]) {
		slog.Debug("onFailed", "triggeredAtUtc", pwd.TriggeredAtUtc, "data", pwd.Data)
	})
}

func (s *SteamcmdInstance) cleanup() {
	s.cmdLock.Lock()
	defer s.cmdLock.Unlock()

	if s.pty != nil {
		s.pty.Close()
	}
	s.pty = nil

	if s.cmd != nil {
		if s.cmd.Process != nil {
			s.cmd.Process.Kill()
			s.cmd.Process.Release()
		}
	}
	s.cmd = nil
}

func (s *SteamcmdInstance) IsRunning() bool {
	return s.running.Load()
}

func (s *SteamcmdInstance) Update(force bool) error {
	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()

	if s.IsRunning() {
		return errors.New("SteamCmdService is busy")
	}

	s.running.Store(true)
	s.onStarted.Trigger()

	if force {
		os.RemoveAll(s.steamCmdDir)
	}

	if !IsSteamCmdInstalled(s.steamCmdDir) {
		if err := downloadSteamCmd(s.steamCmdDir); err != nil {
			s.onFailed.Trigger(err)
			return err
		}
	}

	if err := s.update(); err != nil {
		s.onFailed.Trigger(err)
		return err
	}

	return nil
}

func (s *SteamcmdInstance) update() error {
	var steamCmdShFilePath = filepath.Join(s.steamCmdDir, "steamcmd.sh")

	if err := os.MkdirAll(s.serverDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command(steamCmdShFilePath,
		"+force_install_dir "+s.serverDir,
		"+login anonymous",
		"+app_update 730",
		"validate",
		"+quit",
	)
	f, err := pty.Start(cmd)
	if err != nil {
		s.onFailed.Trigger(err)
		return err
	}

	s.pty = f
	s.cmd = cmd

	go s.checkIfCmdIsRunning(cmd)
	go s.readOutput(f)

	return nil
}

func (s *SteamcmdInstance) checkIfCmdIsRunning(cmd *exec.Cmd) {
	err := cmd.Wait()

	if s.canceled.Load() {
		slog.Debug("steamcmd has been canceled")
		s.onCancelled.Trigger()
	} else if s.lastLine.Load() == "Success! App '730' fully installed." {
		slog.Debug("steamcmd finished")
		s.onFinished.Trigger()
		return
	} else if err != nil {
		slog.Debug("steamcmd exited with error " + err.Error())
		s.onFailed.Trigger(errors.New("steamcmd exited with error " + err.Error()))
		return
	} else {
		slog.Debug("steamcmd exited unexpectedly")
		s.onFailed.Trigger(errors.New("steamcmd exited unexpectedly"))
	}

	slog.Debug("checkIfCmdIsRunning exited")
}

func (s *SteamcmdInstance) readOutput(f *os.File) {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lastLine := scanner.Text()
		lastLine = strings.TrimSpace(lastLine)
		if lastLine == "" {
			continue
		}

		s.onOutput.Trigger(lastLine)
		s.lastLine.Store(lastLine)
	}

	slog.Debug("read output exited. " + scanner.Err().Error())
}

func (s *SteamcmdInstance) Cancel() error {
	if !s.IsRunning() {
		return errors.New("steamcmd is not running")
	}

	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()
	s.cleanup()

	timeout := time.Second * 5
	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout of %q reached", timeout)
		}

		if !s.IsRunning() {
			return nil
		}
	}
}
