package server

import (
	"bufio"
	"cs-server-controller/event"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (s *ServerInstance) Start(sp StartParameters) error {
	spJson, _ := json.Marshal(sp)
	slog.Debug("trying to start server with start parameters " + string(spJson))
	if s.IsRunning() {
		slog.Debug("failed to start server. server is already running")
		return errors.New("server is running")
	}

	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()
	s.running.Store(true)
	s.onServerStarting.Trigger()

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
		s.onServerCrashed.Trigger(err)
		return err
	}
	slog.Debug("server process started")

	s.cmdLock.Lock()
	s.cmd = cmd
	s.writer = inW
	s.cmdLock.Unlock()

	go s.checkIfServerIsRunning(cmd)
	go s.flushServerOutput(inW)
	go s.readServerOutput(outR)

	if err := s.waitForServerToStart(inW, time.Millisecond*30_000); err != nil {
		slog.Debug("failed to wait for server to start. " + err.Error())
		s.onServerCrashed.Trigger(err)
		return err
	}

	s.onServerStarted.Trigger()
	slog.Info("server started")
	return nil
}

func (s *ServerInstance) checkIfServerIsRunning(cmd *exec.Cmd) {
	err := cmd.Wait()

	if s.stop.Load() {
		s.onServerStopped.Trigger()
	} else {
		if err != nil {
			s.onServerCrashed.Trigger(fmt.Errorf("server exited unexpectedly. %v", err))
		} else {
			s.onServerCrashed.Trigger(fmt.Errorf("server exited unexpectedly"))
		}
	}

	slog.Debug("check if server is running exited")
}

func (s *ServerInstance) flushServerOutput(writer *io.PipeWriter) {
	for {
		if !s.IsRunning() {
			break
		}

		s.cmdLock.Lock()
		if writer == nil {
			break
		}

		if _, err := writer.Write([]byte("\n")); err != nil {
			break
		}
		s.cmdLock.Unlock()
		time.Sleep(time.Millisecond * 500)

	}
	slog.Info("server output flush task stopped")
}

func (s *ServerInstance) readServerOutput(reader *io.PipeReader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() && s.IsRunning() {
		out := strings.TrimSpace(scanner.Text())
		if out != "" {
			s.onOutput.Trigger(out)
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Warn("failed to read server output. " + err.Error())
	}

	slog.Debug("readServerOutput exited ")
}

func (s *ServerInstance) waitForServerToStart(write *io.PipeWriter, timeout time.Duration) error {
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
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *ServerInstance) Stop() error {
	if !s.IsRunning() {
		return nil
	}

	slog.Debug("stopping server")

	s.startStopLock.Lock()
	defer s.startStopLock.Unlock()
	s.stop.Store(true)

	if err := s.writeCommand("quit"); err != nil {
		s.cmdLock.Lock()
		if s.cmd == nil {
			s.cmdLock.Unlock()
			return errors.New("quit command failed but process is not running")
		}

		slog.Warn("failed to stop gracefully. Killing process")

		if s.cmd.Process != nil {
			s.cmd.Process.Kill()
			s.cmd.Process.Release()
		}
		s.cmdLock.Unlock()
	}

	const timeout = time.Second * 15
	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout of %q reached", timeout)
		}

		if !s.IsRunning() {
			slog.Info("server stopped")
			return nil
		}

		time.Sleep(time.Second)
	}
}
