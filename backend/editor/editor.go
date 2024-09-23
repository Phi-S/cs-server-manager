package editor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Phi-S/cs-server-manager/files"
	"github.com/Phi-S/cs-server-manager/gvalidator"
)

type FilesToEdit struct {
	Path string `json:"path"`

	// all files with those Extensions in the folder. Requires path to point to an folder
	Extensions []string `json:"extensions"`
}

type Instance struct {
	serverDir   string
	filesToEdit []FilesToEdit
}

func New(editorFilesJsonPath string, serverDir string) (*Instance, error) {
	if err := gvalidator.Instance().Var(editorFilesJsonPath, "filepath"); err != nil {
		return nil, fmt.Errorf("validation editorFilesJsonPath: %w", err)
	}

	if err := gvalidator.Instance().Var(serverDir, "dir"); err != nil {
		return nil, fmt.Errorf("validation baseDir: %w", err)
	}

	var filesToEdit []FilesToEdit

	_, err := os.Stat(editorFilesJsonPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			filesToEdit = []FilesToEdit{
				{Path: "/game/csgo/cfg", Extensions: []string{".cfg"}},
				{Path: "/game/csgo/addons/counterstrikesharp/configs", Extensions: []string{".json"}},
				{Path: "/game/csgo/addons/counterstrikesharp/plugins", Extensions: []string{".json", ".cfg"}},
			}
		} else {
			return nil, fmt.Errorf("os.Stat: %w", err)
		}
	} else {
		jsonRaw, err := os.ReadFile(editorFilesJsonPath)
		if err != nil {
			return nil, fmt.Errorf("read editor json file: %w", err)
		}

		if err := json.Unmarshal(jsonRaw, &filesToEdit); err != nil {
			return nil, fmt.Errorf("json.Unmarshal editor json: %w", err)
		}
	}

	return &Instance{
		serverDir:   serverDir,
		filesToEdit: filesToEdit,
	}, nil
}

func (i *Instance) GetAllEditableFiles() ([]string, error) {
	allFilesInBaseDir, err := files.GetAllFilesInDir(i.serverDir)
	if err != nil {
		return nil, fmt.Errorf("GetAllFilesInDir in baseDir '%v': %w", i.serverDir, err)
	}

	result := make([]string, 0)
	for _, f := range allFilesInBaseDir {
		path := strings.Replace(f.FullPath, i.serverDir, "", 1)
		if i.fileCanBeEdited(path) {
			result = append(result, path)
		}
	}

	return result, nil
}

func (i *Instance) fileCanBeEdited(path string) bool {
	for _, f := range i.filesToEdit {
		if f.Extensions == nil {
			if f.Path == path {
				return true
			}
		} else {
			if f.Path == filepath.Dir(path) {
				for _, ex := range f.Extensions {
					if strings.HasSuffix(path, ex) {
						return true
					}
				}
			}
		}
	}

	return false
}

func (i *Instance) GetFileContent(path string) ([]byte, error) {
	fullPath := filepath.Join(i.serverDir, path)

	if !i.fileCanBeEdited(path) {
		return nil, fmt.Errorf("file can not be edited")
	}

	if err := gvalidator.Instance().Var(fullPath, "filepath"); err != nil {
		return nil, fmt.Errorf("validation full path: %w", err)
	}

	if _, err := os.Stat(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file dose not exist")
		}
		return nil, fmt.Errorf("os.Stat: %w", err)
	}

	result, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return result, nil
}

func (i *Instance) SetFileContent(path string, content []byte) error {
	fullPath := filepath.Join(i.serverDir, path)

	if !i.fileCanBeEdited(path) {
		return fmt.Errorf("file can not be edited")
	}

	if err := gvalidator.Instance().Var(fullPath, "filepath"); err != nil {
		return fmt.Errorf("validation full path: %w", err)
	}

	if _, err := os.Stat(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file dose not exist")
		}
		return fmt.Errorf("os.Stat: %w", err)
	}

	if err := os.WriteFile(fullPath, content, os.ModeAppend.Perm()); err != nil {
		return fmt.Errorf("writeFile: %w", err)
	}

	return nil
}
