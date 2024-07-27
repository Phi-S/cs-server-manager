package json_file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"

	"github.com/asaskevich/govalidator"
)

type json_file[T any] struct {
	path string

	tType reflect.Type
	lock  sync.Mutex
}

func (j *json_file[T]) GetPath() string {
	return j.path
}

func (j *json_file[T]) GetType() reflect.Type {
	return j.tType
}

func (j *json_file[T]) Write(data T) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	j.lock.Lock()
	defer j.lock.Unlock()

	if _, err := os.Stat(j.path); !os.IsNotExist(err) {
		if err := os.Rename(j.path, j.path+".bak"); err != nil {
			return err
		}
	}

	if err := os.WriteFile(j.path, json, 0777); err != nil {
		return err
	}

	return nil
}

func (j *json_file[T]) Read() (*T, error) {
	j.lock.Lock()
	defer j.lock.Unlock()

	content, err := os.ReadFile(j.path)
	if err != nil {
		return nil, err
	}

	data := new(T)
	if err := json.Unmarshal(content, data); err != nil {
		return nil, err
	}

	return data, nil
}

var jsonFileInstances map[string]*json_file[any] = make(map[string]*json_file[any])
var jsonFileInstancesLock sync.Mutex

func Get[T any](pathIn string) (*json_file[T], error) {
	if runtime.GOOS == "windows" && !govalidator.IsWinFilePath(pathIn) {
		return nil, fmt.Errorf("path %q is not a valid win file path", pathIn)
	} else if !govalidator.IsUnixFilePath(pathIn) {
		return nil, fmt.Errorf("path %q is not a valid file path", pathIn)
	}

	pathIn, err := filepath.Abs(pathIn)
	if err != nil {
		return nil, err
	}

	jsonFileInstancesLock.Lock()
	defer jsonFileInstancesLock.Unlock()

	var tType T
	requiredType := reflect.TypeOf(tType)

	if jsonFileInstance, ok := jsonFileInstances[pathIn]; ok {
		instanceType := jsonFileInstance.GetType()
		if instanceType != requiredType {
			return nil, fmt.Errorf("found instance for path %q is of type %q and not of the requested type %q",
				jsonFileInstance.path, instanceType, requiredType)
		}

		return (*json_file[T])(jsonFileInstance), nil
	}

	jsonFileInstance := &json_file[T]{
		path:  pathIn,
		tType: requiredType,
	}

	jsonFileInstances[pathIn] = (*json_file[any])(jsonFileInstance)
	return jsonFileInstance, nil
}
