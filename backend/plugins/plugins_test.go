package plugins_test

import (
	"cs-server-manager/plugins"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	csGamePath := filepath.Join(tempDirPath, "game/")
	err := os.MkdirAll(csGamePath, os.ModePerm)
	if err != nil {
		t.Fatal("os.MkdirAll temp dir", err)
	}

	pluginsJsonPath := filepath.Join(tempDirPath, "plugins.json")
	installedPluginsJsonPath := filepath.Join(tempDirPath, "installed-plugins.json")
	_, err = plugins.New(csGamePath, pluginsJsonPath, installedPluginsJsonPath)
	if err != nil {
		t.Fatal("plugins.New temp dir", err)
	}

	if !t.Failed() {
		defer func() {
			if err := os.RemoveAll(tempDirPath); err != nil {
				t.Log("success but failed to cleanup test dir: ", tempDirPath)
			}
		}()
	}
}

func TestInstallMetamod(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	csGamePath := filepath.Join(tempDirPath, "game/")
	err := os.MkdirAll(csGamePath, os.ModePerm)
	if err != nil {
		t.Fatal("os.MkdirAll temp dir", err)
	}

	pluginsJsonPath := filepath.Join(tempDirPath, "plugins.json")
	installedPluginsJsonPath := filepath.Join(tempDirPath, "installed-plugins.json")
	pluginsInstance, err := plugins.New(csGamePath, pluginsJsonPath, installedPluginsJsonPath)
	if err != nil {
		t.Fatal("plugins.New temp dir", err)
	}

	if err := pluginsInstance.InstallPluginByName("metamod_source", "2.0.0-git1313"); err != nil {
		t.Fatal("InstallPluginByName", err)
	}

	installedPluginsJson, err := os.ReadFile(installedPluginsJsonPath)
	if err != nil {
		t.Fatal("read installed plugins json", err)
	}
	shouldBeJson :=
		`"Name": "metamod_source",
        "Version": "2.0.0-git1313",
        "InstalledAtUtc":`
	if strings.HasPrefix(string(installedPluginsJson), shouldBeJson) {
		t.Fatal("failed to write valid installedJsonFile.json", "is json:", string(installedPluginsJson), "should be json:", shouldBeJson)
	}

	if !t.Failed() {
		defer func() {
			if err := os.RemoveAll(tempDirPath); err != nil {
				t.Log("success but failed to cleanup test dir: ", tempDirPath)
			}
		}()
	}
}

func TestInstallCounterStrikeSharp(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	csGamePath := filepath.Join(tempDirPath, "game/")
	err := os.MkdirAll(csGamePath, os.ModePerm)
	if err != nil {
		t.Fatal("os.MkdirAll temp dir", err)
	}

	pluginsJsonPath := filepath.Join(tempDirPath, "plugins.json")
	installedPluginsJsonPath := filepath.Join(tempDirPath, "installed-plugins.json")
	pluginsInstance, err := plugins.New(csGamePath, pluginsJsonPath, installedPluginsJsonPath)
	if err != nil {
		t.Fatal("plugins.New temp dir", err)
	}

	if err := pluginsInstance.InstallPluginByName("CounterStrikeSharp", "v255"); err != nil {
		t.Fatal("InstallPluginByName", err)
	}

	installedPluginsJson, err := os.ReadFile(installedPluginsJsonPath)
	if err != nil {
		t.Fatal("read installed plugins json", err)
	}
	shouldBeJson :=
		`"Name": "CounterStrikeSharp",
        "Version": "v255",
        "InstalledAtUtc":`
	if strings.HasPrefix(string(installedPluginsJson), shouldBeJson) {
		t.Fatal("failed to write valid installedJsonFile.json", "is json:", string(installedPluginsJson), "should be json:", shouldBeJson)
	}

	if !t.Failed() {
		defer func() {
			if err := os.RemoveAll(tempDirPath); err != nil {
				t.Log("success but failed to cleanup test dir: ", tempDirPath)
			}
		}()
	}
}

func TestInstallCs2PracticeMode(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	csGamePath := filepath.Join(tempDirPath, "game/")
	err := os.MkdirAll(csGamePath, os.ModePerm)
	if err != nil {
		t.Fatal("os.MkdirAll temp dir", err)
	}

	pluginsJsonPath := filepath.Join(tempDirPath, "plugins.json")
	installedPluginsJsonPath := filepath.Join(tempDirPath, "installed-plugins.json")
	pluginsInstance, err := plugins.New(csGamePath, pluginsJsonPath, installedPluginsJsonPath)
	if err != nil {
		t.Fatal("plugins.New temp dir", err)
	}

	if err := pluginsInstance.InstallPluginByName("Cs2PracticeMode", "0.0.14"); err != nil {
		t.Fatal("InstallPluginByName", err)
	}

	installedPluginsJson, err := os.ReadFile(installedPluginsJsonPath)
	if err != nil {
		t.Fatal("read installed plugins json", err)
	}
	shouldBeJson :=
		`"Name": "Cs2PracticeMode",
        "Version": "0.0.14",
        "InstalledAtUtc":`
	if strings.HasPrefix(string(installedPluginsJson), shouldBeJson) {
		t.Fatal("failed to write valid installedJsonFile.json", "is json:", string(installedPluginsJson), "should be json:", shouldBeJson)
	}

	/*
		if !t.Failed() {
			defer func() {
				if err := os.RemoveAll(tempDirPath); err != nil {
					t.Log("success but failed to cleanup test dir: ", tempDirPath)
				}
			}()
		}*/
}
