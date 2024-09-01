package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func Test_getFolderSize_OK(t *testing.T) {
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

	isSize, err := getFolderSize(tempDirPath)
	if err != nil {
		t.Fatal(err)
	}

	if shouldSize != int(isSize) {
		t.Fatalf("Wrong size. Should be: %v | is: %v", shouldSize, isSize)
	}
}
