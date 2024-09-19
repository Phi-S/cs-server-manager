package editor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateExampleEditorFilesJson(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	editorFilesJsonPath := filepath.Join(tempDirPath, "editor-files.json")

	data := []FilesToEdit{
		{Path: "/path/to/file.txt"},
		{Path: "/game/csgo/cfg", Extensions: []string{".cfg"}},
		{Path: "/game/csgo/addons/counterstrikesharp/configs", Extensions: []string{".json"}},
		{Path: "/game/csgo/addons/counterstrikesharp/plugins", Extensions: []string{".json", ".cfg"}},
	}

	json, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(editorFilesJsonPath, json, os.ModePerm); err != nil {
		t.Fatal(err)
	}
}
