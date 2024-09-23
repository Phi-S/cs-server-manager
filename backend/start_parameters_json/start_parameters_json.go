package start_parameters_json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Phi-S/cs-server-manager/event"
	"github.com/Phi-S/cs-server-manager/gvalidator"
	"github.com/Phi-S/cs-server-manager/server"
)

type Instance struct {
	path string
	lock sync.Mutex

	onUpdated event.InstanceWithData[server.StartParameters]
}

func New(path string, defaultIfNotExist server.StartParameters) (*Instance, error) {
	if err := gvalidator.Instance().Var(path, "required,filepath"); err != nil {
		return nil, fmt.Errorf("path validation: %w", err)
	}

	instance := Instance{
		path: path,
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := gvalidator.Instance().Struct(defaultIfNotExist); err != nil {
				return nil, fmt.Errorf("defaultIfNoTExist validation: %w", err)
			}

			if err := instance.Write(defaultIfNotExist); err != nil {
				return nil, fmt.Errorf("instance.Write(defaultIfNotExist): %w", err)
			}
		} else {
			return nil, fmt.Errorf("os.Stat: %w", err)
		}
	}

	return &instance, nil
}

func (i *Instance) Read() (server.StartParameters, error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	var startParameters server.StartParameters

	content, err := os.ReadFile(i.path)
	if err != nil {
		return startParameters, fmt.Errorf("os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(content, &startParameters); err != nil {
		return startParameters, fmt.Errorf("json.Unmarshal: %w", err)
	}

	if err := gvalidator.Instance().Struct(startParameters); err != nil {
		return startParameters, fmt.Errorf("startParameters validation: %w", err)
	}

	return startParameters, nil
}

func (i *Instance) Write(startParameters server.StartParameters) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if err := gvalidator.Instance().Struct(startParameters); err != nil {
		return fmt.Errorf("startParameters validation: %w", err)
	}

	jsonContent, err := json.MarshalIndent(startParameters, "", "    ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	if err := os.WriteFile(i.path, jsonContent, os.ModePerm); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}

	i.onUpdated.Trigger(startParameters)
	return nil
}

func (i *Instance) OnUpdated(handler func(data event.PayloadWithData[server.StartParameters])) {
	i.onUpdated.Register(handler)
}
