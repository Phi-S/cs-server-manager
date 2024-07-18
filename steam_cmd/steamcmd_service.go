package steamcmd_service

import (
	"bufio"
	"cs-server-controller/event"
	"errors"
	"fmt"
	"github.com/creack/pty"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func New(SteamCmdDir string, ServerDir string) SteamcmdService {
	return SteamcmdService{
		busy:        false,
		steamCmdDir: SteamCmdDir,
		serverDir:   ServerDir,
	}
}

type SteamcmdService struct {
	busy                       bool
	steamCmdDir                string
	serverDir                  string
	mutex                      sync.Mutex
	process                    *os.Process
	onUpdateOrInstallStarted   event.Event
	onUpdateOrInstallDone      event.Event
	onUpdateOrInstallCancelled event.Event
	onUpdateOrInstallFailed    event.EventWithData[error]
	onOutput                   event.EventWithData[string]
}

func (s *SteamcmdService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.busy {
		return errors.New("SteamCmdService is busy")
	}

	s.busy = true
	s.onUpdateOrInstallStarted.Trigger()

	if !IsSteamCmdInstalled(s.steamCmdDir) {
		if err := downloadSteamCmd(s.steamCmdDir); err != nil {
			s.onUpdateOrInstallFailed.Trigger(err)
			return err
		}
	}

	if err := s.update(); err != nil {
		s.onUpdateOrInstallFailed.Trigger(err)
		return err
	}

	return nil
}

func (s *SteamcmdService) update() error {
	var steamCmdShFilePath = filepath.Join(s.steamCmdDir, "steamcmd.sh")

	if err := os.MkdirAll(s.serverDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command(steamCmdShFilePath,
		"+force_install_dir "+s.serverDir,
		"+login anonymous",
		"+app_update 740",
		"validate",
		"+quit",
	)
	f, err := pty.Start(cmd)
	if err != nil {
		s.onUpdateOrInstallFailed.Trigger(err)
		return err
	}

	s.process = cmd.Process

	go func() {
		scanner := bufio.NewScanner(f)
		lastLine := ""

		for scanner.Scan() {
			lastLine = scanner.Text()
			s.onOutput.Trigger(lastLine)
		}

		if lastLine == "Success! App '740' fully installed." {
			log.Println("done reading.....")
			s.mutex.Lock()
			s.busy = false
			s.onUpdateOrInstallDone.Trigger()
			s.mutex.Unlock()
		} else {
			err := scanner.Err()
			if err == nil {
				err = errors.New("unknown reason")
			}

			log.Println("Failed to read " + err.Error())
			s.mutex.Lock()
			s.busy = false
			s.onUpdateOrInstallFailed.Trigger(err)
			s.mutex.Unlock()
		}
	}()

	return nil
}

func (s *SteamcmdService) Cancel() error {
	if !s.IsBusy() {
		return errors.New("nothing to cancel")
	}

	s.mutex.Lock()
	if s.process == nil {
		return errors.New("is busy but no process set. this should never happen")
	}

	var err = s.process.Kill()
	if err != nil {
		fmt.Println("failed to kill")
		err = s.process.Release()
		if err != nil {
			fmt.Println("failed to release")
			return err
		}
	}
	defer s.mutex.Unlock()
	return nil
}

func (s *SteamcmdService) IsBusy() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.busy
}

func (s *SteamcmdService) EventOnUpdateOrInstallStarted() *event.Event {
	return &s.onUpdateOrInstallStarted
}

func (s *SteamcmdService) EventOnUpdateOrInstallDone() *event.Event {
	return &s.onUpdateOrInstallDone
}

func (s *SteamcmdService) EventOnUpdateOrInstallCancelled() *event.Event {
	return &s.onUpdateOrInstallCancelled
}

func (s *SteamcmdService) EventOnUpdateOrInstallFailed() *event.EventWithData[error] {
	return &s.onUpdateOrInstallFailed
}

func (s *SteamcmdService) EventOnOutput() *event.EventWithData[string] {
	return &s.onOutput
}
