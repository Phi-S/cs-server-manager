package steamcmd

import (
	download "cs-server-manager/download/unzip"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func IsSteamCmdInstalled(steamCmdPath string) bool {
	if _, err := os.Stat(steamCmdPath); err != nil {
		return false
	}

	if _, err := os.Stat(filepath.Join(steamCmdPath, "steamcmd.sh")); err != nil {
		return false
	}

	if _, err := os.Stat(filepath.Join(steamCmdPath, "linux32")); err != nil {
		return false
	}

	return true
}

func downloadSteamCmd(steamCmdPath string) error {
	if err := os.MkdirAll(steamCmdPath, 0755); err != nil {
		return err
	}

	resp, err := http.Get("https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed with status code " + resp.Status)
	}

	var steamCmdTarGzFilePath = filepath.Join(steamCmdPath, "steamcmd_linux.tar.gz")
	file, err := os.Create(steamCmdTarGzFilePath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	if _, err := download.TarGz(steamCmdTarGzFilePath, steamCmdPath); err != nil {
		return err
	}

	return nil
}
