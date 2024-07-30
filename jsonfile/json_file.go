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

type jsonFile[T any] struct {
    path string

    tType reflect.Type
    lock  sync.Mutex
}

func (j *jsonFile[T]) GetPath() string {
    return j.path
}

func (j *jsonFile[T]) GetType() reflect.Type {
    return j.tType
}

func (j *jsonFile[T]) Write(data T) error {
    dataAsJson, err := json.Marshal(data)
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

func (j *jsonFile[T]) Read() (*T, error) {
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

var jsonFileInstances map[string]*jsonFile[any] = make(map[string]*jsonFile[any])
var jsonFileInstancesLock sync.Mutex

func Get[T any](pathIn string) (*jsonFile[T], error) {
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

        return (*jsonFile[T])(jsonFileInstance), nil
    }

    jsonFileInstance := &jsonFile[T]{
        path:  pathIn,
        tType: requiredType,
    }

    jsonFileInstances[pathIn] = (*jsonFile[any])(jsonFileInstance)
    return jsonFileInstance, nil
}
