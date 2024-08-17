package jfile

import (
	"bytes"
	"cs-server-manager/event"
	"cs-server-manager/gvalidator"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
)

func New[T any](pathIn string, defaultValueIfNoExist T) (*Instance[T], error) {
	if err := gvalidator.Instance().Var(pathIn, "required,filepath"); err != nil {
		return nil, err
	}

	var tType T
	requiredType := reflect.TypeOf(tType)

	jsonFileInstance := &Instance[T]{
		path:  pathIn,
		tType: requiredType,
	}

	_, err := os.Stat(pathIn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := jsonFileInstance.Write(defaultValueIfNoExist); err != nil {
				return nil, fmt.Errorf("jfile.Write defaultValueIfNoExist: %w", err)
			}
		} else {
			return nil, fmt.Errorf("os.Stat: %w", err)
		}
	}

	if _, err := jsonFileInstance.Read(); err != nil {
		return nil, fmt.Errorf("jsonFileInstance.Read first read try: %w", err)
	}

	return jsonFileInstance, nil
}

type Instance[T any] struct {
	path string

	tType reflect.Type
	lock  sync.Mutex

	onUpdated event.InstanceWithData[T]
}

func (j *Instance[T]) Write(data T) error {
	dataAsJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	j.lock.Lock()
	defer j.lock.Unlock()

	if err := os.WriteFile(j.path, dataAsJson, 0777); err != nil {
		return err
	}

	return nil
}

func (j *Instance[T]) Read() (*T, error) {
	j.lock.Lock()
	defer j.lock.Unlock()
	return j.readInternal()
}

func (j *Instance[T]) readInternal() (*T, error) {
	content, err := os.ReadFile(j.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v %w", j.path, err)
	}

	if string(content) == "[]" {
		return new(T), nil
	}

	data := new(T)

	// using decoder so jsons with too many field produce errors
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %v %w", j.path, err)
	}

	if err := gvalidator.Instance().Var(data, "dive"); err != nil {
		return nil, fmt.Errorf("failed to validate file content %v %w", j.path, err)
	}

	return data, nil
}

func (j *Instance[T]) Update(updateFunc func(currentData *T)) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	data, err := j.readInternal()
	if err != nil {
		return fmt.Errorf("failed to read current json file %w", err)
	}

	updateFunc(data)

	dataAsJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal json %w", err)
	}

	if err := os.WriteFile(j.path, dataAsJson, 0777); err != nil {
		return fmt.Errorf("failed to write json file %w", err)
	}

	j.onUpdated.Trigger(*data)
	return nil
}

func (j *Instance[T]) OnUpdated(handler func(data event.PayloadWithData[T])) {
	j.onUpdated.Register(handler)
}
