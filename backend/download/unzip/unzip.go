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

func TarGz(gzFilePath, targetDir string) error {
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zipped file: %w", err)
	}

	defer func() {
		if err := gzFile.Close(); err != nil {
			slog.Warn("UnzipTarGz successful but failed to close source file", "error", err)
		}
	}()

	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}

	defer func() {

		if err := gzReader.Close(); err != nil {
			slog.Warn("UnzipTarGz successful but failed to close reader", "error", err)
		}
	}()

	tarReader := tar.NewReader(gzReader)
	filesToClose := make([]*os.File, 0)
	defer func() {
		for _, f := range filesToClose {
			err := f.Close()
			if err != nil {
				slog.Warn("UnzipTarGz successful but failed to close file result file", "error", err)
			}
		}
	}()

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to download next file %w", err)
		}

		targetFile := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetFile, 0755); err != nil {
				return fmt.Errorf("failed to create new directory for file %v %w", header.Name, err)
			}
		case tar.TypeReg:
			targetDir := filepath.Dir(targetFile)
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create new directory for file %v %w", header.Name, err)
			}

			file, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("failed to create result file for %v %w", header.Name, err)
			}
			filesToClose = append(filesToClose, file)

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: %c in %s", header.Typeflag, header.Name)
		}
	}

	_ = os.Remove(gzFilePath)
	return nil
}

func Zip(zipFilePath, targetDir string) error {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file %w", err)
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

	for _, f := range reader.File {
		filePath := filepath.Join(targetDir, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		filesToClose = append(filesToClose, destinationFile)

		zippedFile, err := f.Open()
		if err != nil {
			return err
		}
		zipFilesToClose = append(zipFilesToClose, zippedFile)

		if _, err := io.Copy(destinationFile, zippedFile); err != nil {
			return err
		}
	}

	return nil
}
