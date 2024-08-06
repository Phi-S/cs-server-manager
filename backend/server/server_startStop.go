package server

import (
	"bufio"
	"cs-server-manager/event"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (s *Instance) copySteamclient() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir %w", err)
	}

	steamClientName := "steamclient.so"
	steamClientSrcPath := filepath.Join(s.steamcmdDir, "linux64", steamClientName)
	steamClientDestPath := filepath.Join(homeDir, ".steam", "sdk64", steamClientName)

	if _, err := os.Stat(steamClientSrcPath); err != nil {
		return fmt.Errorf("steamclient source not found %w", err)
	}

	src, err := os.ReadFile(steamClientSrcPath)
	if err != nil {
		return fmt.Errorf("failed to read steamclient source %w", err)
	}

	_, err = os.Stat(steamClientDestPath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		fmt.Println("err: " + err.Error())
		return fmt.Errorf("unexpected error while checking if steamclient already exists %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(steamClientDestPath), 0770); err != nil {
		return fmt.Errorf("failed to create directories for steamclient %w", err)
	}

	if err := os.WriteFile(steamClientDestPath, src, 0770); err != nil {
		return fmt.Errorf("failed to write steamclient file %w", err)
	}

	return nil
}

func (s *Instance) waitForServerToExit(onExitedChannel chan error) {
	err := <-onExitedChannel

	if s.stop.Load() {
		s.onStopped.Trigger()
	} else if err != nil {
		s.onCrashed.Trigger(fmt.Errorf("server exited unexpectedly. %v", err))
	} else {
		s.onCrashed.Trigger(fmt.Errorf("server exited unexpectedly"))
	}

	s.Close()
	s.running.Store(false)
	s.stop.Store(false)
	s.started.Store(false)
	slog.Debug("server exited and all resources released")
}

func (s *Instance) Start(sp StartParameters) error {
	if s.IsRunning() {
		return errors.New("server is running")
	}

	slog.Debug("trying to start server with start parameters", "start-parameters", sp)

	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()

	onExitedChan := make(chan error)
	go s.waitForServerToExit(onExitedChan)
	s.running.Store(true)
	s.onStarting.Trigger()

	if err := s.copySteamclient(); err != nil {
		onExitedChan <- err
		return fmt.Errorf("failed to copy steamclient %w", err)
	}

	password := strings.TrimSpace(sp.Password)
	if len([]rune(password)) == 0 {
		password = ""
	} else {
		password = " +sv_password " + password
	}

	loginToken := strings.TrimSpace(sp.SteamLoginToken)
	if len([]rune(loginToken)) == 0 {
		loginToken = ""
	} else {
		loginToken = " +sv_setsteamaccount " + password
	}
	//TODO: use sp.Additional

	cs2ShPath := filepath.Join(s.serverDir, "game", "bin", "linuxsteamrt64", "cs2")
	cmd := exec.Command(cs2ShPath,
		"-dedicated",
		"-console",
		fmt.Sprintf("-port %s", s.port),
		fmt.Sprintf("+hostname '%s'", sp.Hostname),
		fmt.Sprintf("-maxplayers %d", sp.MaxPlayers),
		fmt.Sprintf("+map %s", sp.StartMap),
		password,
		loginToken,
	)

	slog.Debug("start command: " + strings.Join(cmd.Args, " "))

	outR, outW := io.Pipe()
	inR, inW := io.Pipe()
	cmd.Stdout = outW
	cmd.Stderr = outW
	cmd.Stdin = inR

	if err := cmd.Start(); err != nil {
		slog.Debug("failed to start server process. " + err.Error())
		onExitedChan <- err
		return err
	}
	slog.Debug("server process started")

	s.cmd = cmd
	s.writer = inW
	s.reader = outR

	go func() {
		err := cmd.Wait()
		onExitedChan <- err
	}()

	go s.flushServerOutput(inW)
	go s.readServerOutput(outR)

	if err := s.waitForServerToStart(inW, time.Millisecond*30_000); err != nil {
		slog.Debug("failed to wait for server to start. " + err.Error())
		onExitedChan <- err
		return err
	}

	s.onStarted.Trigger(sp)
	slog.Info("server started")
	return nil
}

func (s *Instance) flushServerOutput(writer *io.PipeWriter) {
	for {
		if !s.IsRunning() {
			break
		}

		if writer == nil {
			break
		}

		if _, err := writer.Write([]byte("\n")); err != nil {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	slog.Info("server output flush task stopped")
}

func (s *Instance) readServerOutput(reader *io.PipeReader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		out := strings.TrimSpace(scanner.Text())
		if out == "" {
			continue
		}

		s.onOutput.Trigger(out)
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("failed to read server output. " + err.Error())
	}

	slog.Debug("readServerOutput exited ")
}

func (s *Instance) waitForServerToStart(write *io.PipeWriter, timeout time.Duration) error {
	slog.Debug("waiting for server to start")

	const startMessage = "#####_SERVER_STARTED"

	hostActivated := false
	serverStarted := false
	handlerUuid := s.onOutput.Register(func(pwd event.PayloadWithData[string]) {
		if !hostActivated {
			if strings.HasPrefix(strings.TrimSpace(pwd.Data), "Host activate: Loading") {
				write.Write([]byte("say " + startMessage + "\n"))
				hostActivated = true
				slog.Debug("host activated")
			}
		} else {
			if pwd.Data == "[All Chat][Console (0)]: "+startMessage {
				serverStarted = true
			}
		}
	})

	defer s.onOutput.Deregister(handlerUuid)

	startTime := time.Now()
	for {
		if time.Since(startTime) >= timeout {
			return errors.New("timeout reached")
		}

		if serverStarted {
			s.started.Store(true)
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *Instance) Stop() error {
	if !s.IsRunning() {
		return errors.New("server is not running")
	}

	slog.Debug("stopping server")

	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()
	s.stop.Store(true)

	if err := s.writeCommand("quit"); err != nil {
		if s.cmd == nil {
			return errors.New("quit command failed but process is not running")
		}

		slog.Warn("failed to stop gracefully. Killing process")

		if s.cmd.Process != nil {
			_ = s.cmd.Process.Kill()
			_ = s.cmd.Process.Release()
		}
	}

	const gracefulStopTimeout = time.Second * 10
	const stopTimeout = time.Second * 20
	startTime := time.Now()
	for {
		if time.Since(startTime) > gracefulStopTimeout {
			slog.Warn("failed to stop gracefully. Timeout reached. Killing process")
			if s.cmd.Process != nil {
				_ = s.cmd.Process.Kill()
				_ = s.cmd.Process.Release()
			} else {
				return fmt.Errorf("failed to force stop server. Process is nil")
			}
		}

		if time.Since(startTime) > stopTimeout {
			return fmt.Errorf("timeout of %q reached", stopTimeout)
		}

		if !s.IsRunning() {
			slog.Info("server stopped")
			return nil
		}

		time.Sleep(time.Second)
	}
}
