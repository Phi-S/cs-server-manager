package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort    string
	CsPort      string
	DataDir     string
	LogDir      string
	ServerDir   string
	SteamcmdDir string
}

func GetRequiredValueFromEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return "", fmt.Errorf("failed to get %q from environment", key)
	}

	return value, nil
}

func GetOptionalValueFromEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return defaultValue
	}

	return value
}

func GetConfig() (Config, error) {
	const envFile = ".env"
	if err := godotenv.Load(envFile); err != nil {
		return Config{}, fmt.Errorf("failed to load %q file", envFile)
	}

	const httpPortKey = "HTTP_PORT"
	httpPort, err := GetRequiredValueFromEnv(httpPortKey)
	if err != nil {
		return Config{}, err
	}

	if !govalidator.IsPort(httpPort) {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid PORT", httpPortKey, httpPort)
	}
	slog.Info("config", httpPortKey, httpPort)

	const csPortKey = "CS_PORT"
	csPort, err := GetRequiredValueFromEnv(csPortKey)
	if err != nil {
		return Config{}, err
	}
	if !govalidator.IsPort(csPort) {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid PORT", csPortKey, csPort)
	}
	slog.Info("config", csPortKey, csPort)

	const dataDirKey = "DATA_DIR"
	dataDir, err := GetRequiredValueFromEnv(dataDirKey)
	if err != nil {
		return Config{}, err
	}
	if ok, _ := govalidator.IsFilePath(dataDir); !ok {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid filepath", dataDirKey, dataDir)
	}
	slog.Info("config", dataDirKey, dataDir)

	const logDirKey = "LOG_DIR"
	logDir := GetOptionalValueFromEnv(logDirKey, filepath.Join(dataDir, "logs"))
	slog.Info("config", logDirKey, logDir)

	const serverDirKey = "SERVER_DIR"
	serverDir := GetOptionalValueFromEnv(serverDirKey, filepath.Join(dataDir, "server"))
	slog.Info("config", serverDirKey, serverDir)

	const steamcmdDirKey = "STEAMCMD_DIR"
	steamcmdDir := GetOptionalValueFromEnv(steamcmdDirKey, filepath.Join(dataDir, "steamcmd"))
	slog.Info("config", steamcmdDirKey, steamcmdDir)

	return Config{
		httpPort,
		csPort,
		dataDir,
		logDir,
		serverDir,
		steamcmdDir,
	}, nil
}
