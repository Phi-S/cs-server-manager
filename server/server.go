package server

import (
	"cs-server-controller/event"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asaskevich/govalidator"
)

var instanceCreated = false

type ServerInstance struct {
	serverDir string
	port      string

	running atomic.Bool
	stop    atomic.Bool

	cmd    *exec.Cmd
	writer *io.PipeWriter

	startStopLock sync.Mutex
	commandLock   sync.Mutex

	onOutput         event.EventWithData[string]
	onServerStarting event.Event
	onServerStarted  event.Event
	onServerStopped  event.Event
	onServerCrashed  event.EventWithData[error]
}

func NewInstance(serverDir, port string, enableEventLogging bool) (*ServerInstance, error) {
	if instanceCreated {
		return nil, errors.New("another instance already exists. Only one instance should be used throughout the program")
	}

	ok, _ := govalidator.IsFilePath(serverDir)
	if !ok {
		errorMsg := fmt.Sprintf("server dir %q is not a valid filepath", serverDir)
		return nil, errors.New(errorMsg)
	}

	isPort := govalidator.IsPort(port)
	if !isPort {
		errorMsg := fmt.Sprintf("port %q is not a valid filepath", port)
		return nil, errors.New(errorMsg)
	}

	i := ServerInstance{
		serverDir: serverDir,
		port:      port,
	}

	if enableEventLogging {
		i.enableEventLogging()
	}

	i.onServerCrashed.Register(func(pwd event.PayloadWithData[error]) {
		i.Close()
		i.running.Store(false)
		i.stop.Store(false)
	})

	i.onServerStopped.Register(func(dp event.DefaultPayload) {
		i.Close()
		i.running.Store(false)
		i.stop.Store(false)
	})

	instanceCreated = true
	return &i, nil
}

func (s *ServerInstance) IsRunning() bool {
	return s.running.Load()
}

func (s *ServerInstance) Close() {
	if s.cmd != nil {
		if s.cmd.Process != nil {
			s.cmd.Process.Kill()
			s.cmd.Process.Release()
		}
		s.cmd = nil
	}

	// TODO: add in and out reader/writer to struct so they can get closed??
	if s.writer != nil {
		s.writer.Close()
		s.writer = nil
	}
}

func (s *ServerInstance) enableEventLogging() {
	s.onOutput.Register(func(pwd event.PayloadWithData[string]) {
		//slog.Debug("EVENT TRIGGERED", "event", pwd.EventName, "triggeredAtUtc", pwd.TriggeredAtUtc, "data", pwd.Data)
		fmt.Println(pwd.TriggeredAtUtc.Format(time.RFC3339Nano) + " | SERVER: " + pwd.Data)
	})

	s.onServerStarting.Register(func(dp event.DefaultPayload) {
		slog.Debug("onServerStarting", "triggeredAtUtc", dp.TriggeredAtUtc)
	})

	s.onServerStarted.Register(func(dp event.DefaultPayload) {
		slog.Debug("onServerStarted", "triggeredAtUtc", dp.TriggeredAtUtc)
	})

	s.onServerCrashed.Register(func(pwd event.PayloadWithData[error]) {
		slog.Debug("onServerCrashed", "triggeredAtUtc", pwd.TriggeredAtUtc, "data", pwd.Data)
	})

	s.onServerStopped.Register(func(dp event.DefaultPayload) {
		slog.Debug("onServerStopped", "triggeredAtUtc", dp.TriggeredAtUtc)
	})
}
