package steamcmd_service

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
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

	if err := unzip(steamCmdTarGzFilePath, steamCmdPath); err != nil {
		return err
	}

	return nil
}

func unzip(gzFilePath, targetDir string) error {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		targetFile := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFile, 0755); err != nil {
				return err
			}
		case tar.TypeReg:

			targetDir := filepath.Dir(targetFile)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %c in %s", header.Typeflag, header.Name)
		}
	}

	os.Remove(gzFilePath)
	return nil
}
