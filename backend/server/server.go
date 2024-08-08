package server

import (
	"cs-server-manager/event"
	"cs-server-manager/gvalidator"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"sync/atomic"
)

var instanceCreated = false

type Instance struct {
	steamcmdDir string
	serverDir   string
	port        string

	running atomic.Bool
	started atomic.Bool
	stop    atomic.Bool

	cmd    *exec.Cmd
	writer *io.PipeWriter
	reader *io.PipeReader

	startStopLock sync.Mutex
	commandLock   sync.Mutex

	onOutput   event.InstanceWithData[string]
	onStarting event.Instance
	onStarted  event.InstanceWithData[StartParameters]
	onStopped  event.Instance
	onCrashed  event.InstanceWithData[error]
}

func NewInstance(serverDir, port, steamcmdDir string) (*Instance, error) {
	if instanceCreated {
		return nil, errors.New("another instance already exists. Only one instance should be used throughout the program")
	}

	if err := gvalidator.Instance.Var(serverDir, "required,dir"); err != nil {
		return nil, fmt.Errorf("server dir %v is not a valid filepath %w", serverDir, err)
	}

	if err := gvalidator.Instance.Var(port, "required,port"); err != nil {
		return nil, fmt.Errorf("port %q is not a valid filepath %w", port, err)
	}

	if err := gvalidator.Instance.Var(steamcmdDir, "required,dir"); err != nil {
		return nil, fmt.Errorf("steamcmd dir %v is not a valid filepath %w", steamcmdDir, err)
	}

	i := Instance{
		steamcmdDir: steamcmdDir,
		serverDir:   serverDir,
		port:        port,
	}

	i.onCrashed.Register(func(pwd event.PayloadWithData[error]) {
		i.Close()
		i.running.Store(false)
		i.stop.Store(false)
	})

	i.onStopped.Register(func(dp event.DefaultPayload) {
		i.Close()
		i.running.Store(false)
		i.stop.Store(false)
	})

	instanceCreated = true
	return &i, nil
}

func (s *Instance) IsRunning() bool {
	return s.running.Load()
}

func (s *Instance) Close() {
	if s.cmd != nil {
		if s.cmd.Process != nil {
			_ = s.cmd.Process.Kill()
			_ = s.cmd.Process.Release()
		}
		s.cmd = nil
	}

	// TODO: add in and out reader/writer to struct so they can get closed??
	if s.writer != nil {
		_ = s.writer.Close()
		s.writer = nil
	}

	if s.reader != nil {
		_ = s.reader.Close()
		s.reader = nil
	}
}
