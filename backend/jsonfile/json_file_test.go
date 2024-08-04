package json_file_test

import (
	json_file "cs-server-manager/jsonfile"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type TestJsonFile struct {
	Name  string `json:"name" validate:"required"`
	Count int    `json:"count" validate:"required"`
	Bool  bool   `json:"bool" validate:"required"`
}

func TestNew_OK_CreateNewJsonFile(t *testing.T) {
	testDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Error(err)
	}

	testJsonFilePath := filepath.Join(testDir, "test-json.json")
	t.Logf("testing for json file %q", testJsonFilePath)
	_, err = json_file.New[TestJsonFile](testJsonFilePath, TestJsonFile{
		Name:  "name_test",
		Count: 5,
		Bool:  true,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNew_FAIL_NotAValidFilePath(t *testing.T) {
	testDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Error(err)
	}

	testJsonFilePath := filepath.Join(testDir, "test-json.json/##")
	t.Logf("testing for json file %q", testJsonFilePath)
	_, err = json_file.New[TestJsonFile](testJsonFilePath, TestJsonFile{
		Name:  "name_test",
		Count: 5,
		Bool:  true,
	})
	if err == nil {
		t.Error("error expected but nil returned")
	} else if err.Error() != fmt.Sprintf("open %v: no such file or directory", testJsonFilePath) {
		t.Error(err)
	}
}

func TestNew_FAIL_MalformedJsonTooManyField(t *testing.T) {
	testDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Error(err)
	}

	testJsonFilePath := filepath.Join(testDir, "test-json.json")
	t.Logf("testing for json file %q", testJsonFilePath)
	if err := os.WriteFile(
		testJsonFilePath,
		[]byte("{\"name\":\"709\",\"count\":5,\"bool\":true,\"test\":1}"), 0777,
	); err != nil {
		t.Error(err)
	}

	_, err = json_file.New[TestJsonFile](testJsonFilePath, TestJsonFile{
		Name:  "name_test",
		Count: 5,
		Bool:  true,
	})

	if err == nil {
		t.Error("error expected but nil returned")
	} else if err.Error() != "json: unknown field \"test\"" {
		t.Error(err)
	} else {
		t.Logf("expected error received %q", err)
	}
}

func TestNew_FAIL_MalformedJsonMissingField(t *testing.T) {
	testDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Error(err)
	}

	testJsonFilePath := filepath.Join(testDir, "test-json.json")
	t.Logf("testing for json file %q", testJsonFilePath)
	if err := os.WriteFile(
		testJsonFilePath,
		[]byte("{\"name\":\"709\",\"count\":5}"), 0777,
	); err != nil {
		t.Error(err)
	}

	_, err = json_file.New[TestJsonFile](testJsonFilePath, TestJsonFile{
		Name:  "name_test",
		Count: 5,
		Bool:  true,
	})
	if err == nil {
		t.Error("error expected but nil returned")
	} else if err.Error() != "Key: 'TestJsonFile.Bool' Error:Field validation for 'Bool' failed on the 'required' tag" {
		t.Error(err)
	} else {
		t.Logf("expected error received %q", err)
	}
}

func TestNew_OK_ReadExistingJsonFile(t *testing.T) {
	testDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Error(err)
	}

	testJsonFilePath := filepath.Join(testDir, "test-json.json")
	t.Logf("testing for json file %q", testJsonFilePath)
	if err := os.WriteFile(
		testJsonFilePath,
		[]byte("{\"name\":\"709\",\"count\":5,\"bool\":true}"), 0777,
	); err != nil {
		t.Error(err)
	}

	_, err = json_file.New[TestJsonFile](testJsonFilePath, TestJsonFile{
		Name:  "name_test",
		Count: 5,
		Bool:  true,
	})
	if err != nil {
		t.Error(err)
	}
}
