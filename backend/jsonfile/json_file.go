package json_file

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

func New[T any](pathIn string, defaultValueIfNoExist T) (*JsonFile[T], error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	if err := v.Var(pathIn, "required,filepath"); err != nil {
		return nil, err
	}

	var tType T
	requiredType := reflect.TypeOf(tType)

	jsonFileInstance := &JsonFile[T]{
		path:     pathIn,
		tType:    requiredType,
		validate: v,
	}

	_, err := os.Stat(pathIn)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := jsonFileInstance.Write(defaultValueIfNoExist); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if _, err := jsonFileInstance.Read(); err != nil {
		return nil, err
	}

	return jsonFileInstance, nil
}

type JsonFile[T any] struct {
	path string

	tType    reflect.Type
	lock     sync.Mutex
	validate *validator.Validate
}

func (j *JsonFile[T]) GetPath() string {
	return j.path
}

func (j *JsonFile[T]) GetType() reflect.Type {
	return j.tType
}

func (j *JsonFile[T]) Write(data T) error {
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

func (j *JsonFile[T]) Read() (*T, error) {
	j.lock.Lock()
	defer j.lock.Unlock()

	content, err := os.ReadFile(j.path)
	if err != nil {
		return nil, err
	}

	data := new(T)

	// using decoder so jsons with too many field produce errors
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	val := validator.New(validator.WithRequiredStructEnabled())
	if err := val.Struct(data); err != nil {
		return nil, err
	}

	return data, nil
}
