package server

import (
	"cs-server-controller/event"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func (s *ServerInstance) writeCommand(command string) error {
	if !s.IsRunning() {
		return errors.New("server not running")
	}

	if s.writer == nil {
		slog.Error("server is marked as running but writer is not set. this should never happen")
		return errors.New("write is nil")
	}

	commandAsBytes := []byte(command + "\n")
	if _, err := s.writer.Write(commandAsBytes); err != nil {
		slog.Warn("failed to write. " + err.Error())
		return err
	}

	slog.Debug("command send \"" + command + "\"")
	return nil
}

func (s *ServerInstance) SendCommand(command string) (string, error) {
	const startPrefix = "#####_START_"
	const endPrefix = "#####_END_"
	uuid := uuid.New().String()

	start := fmt.Sprintf("%s%s", startPrefix, uuid)
	end := fmt.Sprintf("%s%s", endPrefix, uuid)
	finalCommand := fmt.Sprintf("echo %s\n%s\necho %s", start, command, end)

	s.commandLock.Lock()
	defer s.commandLock.Unlock()

	lock := sync.Mutex{}
	var output = strings.Builder{}
	captureOutput := false
	commandFinished := false
	handlerUuid := s.onOutput.Register(func(pwd event.PayloadWithData[string]) {
		lock.Lock()
		defer lock.Unlock()

		if commandFinished {
			return
		}

		if strings.TrimSpace(pwd.Data) == start {
			captureOutput = true
			return
		}

		if captureOutput && strings.TrimSpace(pwd.Data) == end {
			captureOutput = false
			commandFinished = true
			return
		}

		if captureOutput {
			output.WriteString(pwd.Data + "\n")
		}
	})

	defer s.onOutput.Deregister(handlerUuid)

	if err := s.writeCommand(finalCommand); err != nil {
		return "", err
	}

	timeout := time.Second * 10
	startTime := time.Now()
	for {
		lock.Lock()
		if commandFinished {
			return output.String(), nil
		}
		lock.Unlock()

		if time.Since(startTime) >= timeout {
			return "", errors.New("timeout reached")
		}
		time.Sleep(time.Millisecond * 50)
	}
}
