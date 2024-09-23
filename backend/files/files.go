package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Phi-S/cs-server-manager/gvalidator"
)

type Entry struct {
	FullPath  string
	SizeBytes int64
	ParentDir string
}

func GetAllFilesInDir(dir string) ([]Entry, error) {
	if err := gvalidator.Instance().Var(dir, "dir"); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("readDir: %w", err)
	}

	results := make([]Entry, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		path := filepath.Join(dir, dirEntry.Name())
		if dirEntry.IsDir() {
			filesInSubDir, err := GetAllFilesInDir(path)
			if err != nil {
				return nil, fmt.Errorf("get files in subdir '%v': %w", path, err)
			}

			results = append(results, filesInSubDir...)
		} else {
			info, err := dirEntry.Info()
			if err != nil {
				return nil, fmt.Errorf("dir entry info '%v': %w", path, err)
			}
			results = append(results, Entry{FullPath: path, SizeBytes: info.Size(), ParentDir: dir})
		}
	}

	return results, nil
}

func GetDirSize(dir string) (int64, error) {
	var size int64
	err := filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get size of directory '%v': %w", dir, err)
	}
	return size, nil
}
