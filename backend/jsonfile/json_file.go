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

type JsonFile[T any] struct {
    path string

    tType reflect.Type
    lock  sync.Mutex
}

func (j *JsonFile[T]) GetPath() string {
    return j.path
}

func (j *JsonFile[T]) GetType() reflect.Type {
    return j.tType
}

func (j *JsonFile[T]) Write(data T) error {
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

func (j *JsonFile[T]) Read() (*T, error) {
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

func Get[T any](pathIn string, defaultValueIfNoExist T) (*JsonFile[T], error) {
    if runtime.GOOS == "windows" && !govalidator.IsWinFilePath(pathIn) {
        return nil, fmt.Errorf("path %q is not a valid win file path", pathIn)
    } else if !govalidator.IsUnixFilePath(pathIn) {
        return nil, fmt.Errorf("path %q is not a valid file path", pathIn)
    }

    pathIn, err := filepath.Abs(pathIn)
    if err != nil {
        return nil, err
    }

    var tType T
    requiredType := reflect.TypeOf(tType)

    jsonFileInstance := &JsonFile[T]{
        path:  pathIn,
        tType: requiredType,
    }

    _, err = os.Stat(pathIn)
    if err != nil {
        if os.IsNotExist(err) {
            if err := jsonFileInstance.Write(defaultValueIfNoExist); err != nil {
                return nil, err
            }
        }

        return nil, err
    }

    if _, err := jsonFileInstance.Read(); err != nil {
        return nil, err
    }

    return jsonFileInstance, nil
}
