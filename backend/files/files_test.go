package files

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestGetAllFilesInDir_oneFolderNoFile(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	folder := uuid.NewString()
	folderPath := filepath.Join(tempDirPath, folder)
	if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	files, err := GetAllFilesInDir(tempDirPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 0 {
		t.Fatalf("expected 0 files but found %d", len(files))
	}
}

func TestGetAllFilesInDir_oneFileInOneFolder(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	folderPath := filepath.Join(tempDirPath, uuid.NewString())
	if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	fileInFolder := filepath.Join(folderPath, uuid.NewString())
	if err := os.WriteFile(fileInFolder, []byte("test file content"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	files, err := GetAllFilesInDir(tempDirPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("should have found 1 file but got %d", len(files))
	}
}

func Test_GetDirSize_OK(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	shouldSize := 1024 * 1024 * 1024
	sampleFile := make([]byte, shouldSize)
	rand.Read(sampleFile)

	sampleFilePath := filepath.Join(tempDirPath, "sample")
	err := os.WriteFile(sampleFilePath, sampleFile, os.ModeAppend)
	if err != nil {
		t.Fatal(err)
	}

	isSize, err := GetDirSize(tempDirPath)
	if err != nil {
		t.Fatal(err)
	}

	if shouldSize != int(isSize) {
		t.Fatalf("Wrong size. Should be: %v | is: %v", shouldSize, isSize)
	}
}
