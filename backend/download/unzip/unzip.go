package unzip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func TarGz(gzFilePath, targetDir string) ([]string, error) {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zipped file: %w", err)
	}

	defer func() {
		if err := gzFile.Close(); err != nil {
			slog.Warn("UnzipTarGz successful but failed to close source file", "path", gzFilePath, "error", err)
		}
	}()

	gzipReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip gzipReader: %w", err)
	}

	defer func() {
		if err := gzipReader.Close(); err != nil {
			slog.Warn("UnzipTarGz successful but failed to close gzipReader", "error", err)
		}
	}()

	tarReader := tar.NewReader(gzipReader)

	filesToClose := make([]*os.File, 0)
	defer func() {
		for _, f := range filesToClose {
			err := f.Close()
			if err != nil {
				slog.Warn("UnzipTarGz successful but failed to close file result file", "error", err)
			}
		}
	}()

	extractedFiles := make([]string, 0)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to extract targetFile %w", err)
		}

		targetFilePath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFilePath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create new directory for targetFile %v %w", header.Name, err)
			}
		case tar.TypeReg:
			targetDir := filepath.Dir(targetFilePath)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create new directory for targetFile %v %w", header.Name, err)
			}

			targetFile, err := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create result targetFile for %v %w", header.Name, err)
			}
			filesToClose = append(filesToClose, targetFile)

			if _, err := io.Copy(targetFile, tarReader); err != nil {
				return nil, fmt.Errorf("failed to copy zipped file '%v' to destination '%v' %w", header.Name, targetFilePath, err)
			}
			extractedFiles = append(extractedFiles, targetFilePath)
		default:
			return nil, fmt.Errorf("unsupported type: %c in %s", header.Typeflag, header.Name)
		}
	}

	_ = os.Remove(gzFilePath)
	return extractedFiles, nil
}

func Zip(zipFilePath, targetDir string) ([]string, error) {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file %w", err)
	}
	defer func() {
		if err := reader.Close(); err != nil {
			slog.Warn("Unzip zip successful but failed to close source reader", "error", err)
		}
	}()

	filesToClose := make([]*os.File, 0)
	zipFilesToClose := make([]io.ReadCloser, 0)
	defer func() {
		for _, f := range filesToClose {
			if err := f.Close(); err != nil {
				slog.Warn("UnzipTarGz successful but failed to close file result file", "error", err)
			}
		}

		for _, closer := range zipFilesToClose {
			if err := closer.Close(); err != nil {
				slog.Warn("UnzipTarGz successful but failed to close reader result file", "error", err)
			}
		}
	}()

	extractedFiles := make([]string, 0)

	for _, f := range reader.File {
		filePath := filepath.Join(targetDir, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return nil, fmt.Errorf("invalid file path: %s", filePath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create directory '%v' %w", filePath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory for file '%v' %w", filePath, err)
		}

		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, fmt.Errorf("failed to open destination file '%v' %w", filePath, err)
		}
		filesToClose = append(filesToClose, destinationFile)

		zippedFile, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open zipped file %w", err)
		}
		zipFilesToClose = append(zipFilesToClose, zippedFile)

		if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return nil, fmt.Errorf("failed to copy unzipped file '%v' to destination '%v' %w",
				f.Name, destinationFile, err)
		}
		extractedFiles = append(extractedFiles, filePath)
	}

	if len(extractedFiles) == 0 {
		return nil, fmt.Errorf("no files extracted")
	}

	return extractedFiles, nil
}
