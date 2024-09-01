package download

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Download(url string, destDir string) (FilePath string, error error) {
	split := strings.Split(url, "/")
	fileName := split[len(split)-1]
	destFilePath := filepath.Join(destDir, fileName)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.Warn("resp.Body.Close: finished successfully but failed to close body after download %w", "error", err)
		}
	}()

	out, err := os.Create(destFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			slog.Warn("out.Close: finished successfully but failed to close output file after download %w", "error", err)
		}
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to create result file: %w", err)
	}

	return destFilePath, nil
}
