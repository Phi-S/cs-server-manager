package start_parameters_json_test

import (
	"cs-server-manager/server"
	"cs-server-manager/start_parameters_json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"testing"
)

func TestNew_createNewFile(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_start_parameters_json_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	jsonPath := filepath.Join(tempDirPath, "test.json")
	_, err := start_parameters_json.New(jsonPath, *server.DefaultStartParameters())
	if err != nil {
		t.Fatal(err)
	}

	jsonFileContent, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(jsonFileContent) != `{
    "hostname": "cs server",
    "password": "",
    "start_map": "de_mirage",
    "max_players": 10,
    "steam_login_token": "",
    "additional": []
}` {
		t.Fatal("content dose not match")
	}
}

func TestNew_existingFile(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_start_parameters_json_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	jsonPath := filepath.Join(tempDirPath, "test.json")

	testJson := `{
    "hostname": "test123",
    "password": "",
    "start_map": "de_test",
    "max_players": 111,
    "steam_login_token": "",
    "additional": []
}`

	if err := os.WriteFile(jsonPath, []byte(testJson), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	_, err := start_parameters_json.New(jsonPath, *server.DefaultStartParameters())
	if err != nil {
		t.Fatal(err)
	}

	jsonFileContent, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(jsonFileContent) != testJson {
		t.Fatal("content dose not match")
	}
}

func TestInstance_Read_existingFile(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_start_parameters_json_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	jsonPath := filepath.Join(tempDirPath, "test.json")

	testJson := `{
    "hostname": "test123",
    "password": "",
    "start_map": "de_test",
    "max_players": 111,
    "steam_login_token": "",
    "additional": []
}`

	if err := os.WriteFile(jsonPath, []byte(testJson), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	instance, err := start_parameters_json.New(jsonPath, *server.DefaultStartParameters())
	if err != nil {
		t.Fatal(err)
	}

	startParameters, err := instance.Read()
	if err != nil {
		t.Fatal(err)
	}

	if startParameters.Hostname != "test123" {
		t.Fatal("hostname dose not match")
	}

	if startParameters.Password != "" {
		t.Fatal("password dose not match")
	}

	if startParameters.StartMap != "de_test" {
		t.Fatal("start map dose not match")
	}

	if startParameters.MaxPlayers != 111 {
		t.Fatal("max players dose not match")
	}
}

func TestInstance_Write(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_start_parameters_json_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	jsonPath := filepath.Join(tempDirPath, "test.json")
	instance, err := start_parameters_json.New(jsonPath, *server.DefaultStartParameters())
	if err != nil {
		t.Fatal(err)
	}

	if err := instance.Write(server.StartParameters{
		Hostname:        "test_123",
		Password:        "",
		StartMap:        "de_test123",
		MaxPlayers:      14,
		SteamLoginToken: "",
		Additional:      make([]string, 0),
	}); err != nil {
		t.Fatal(err)
	}

	jsonFileContent, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	if string(jsonFileContent) != `{
    "hostname": "test_123",
    "password": "",
    "start_map": "de_test123",
    "max_players": 14,
    "steam_login_token": "",
    "additional": []
}` {
		t.Fatal("content dose not match")
	}

}
