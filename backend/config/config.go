package config

import (
	globalvalidator "cs-server-manager/global_validator"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

func GetRequiredValueFromEnvAndValidate(key string, validationString string) (string, error) {
	value := os.Getenv(key)
	value = strings.TrimSpace(value)

	if err := globalvalidator.Instance.Var(value, validationString); err != nil {
		return "", fmt.Errorf("validation of %q with the validation string %q and value %q returned error: %q", key, validationString, value, err)
	}

	return value, nil
}

func GetEnvWithDefaultValue(key string, validationString string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return defaultValue
	}

	if err := globalvalidator.Instance.Var(value, validationString); err != nil {
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
	httpPort, err := GetRequiredValueFromEnvAndValidate(httpPortKey, "required,port")
	if err != nil {
		return Config{}, err
	}
	slog.Info("config", httpPortKey, httpPort)

	const csPortKey = "CS_PORT"
	csPort, err := GetRequiredValueFromEnvAndValidate(csPortKey, "required,port")
	if err != nil {
		return Config{}, err
	}
	slog.Info("config", csPortKey, csPort)

	const dataDirKey = "DATA_DIR"
	dataDir, err := GetRequiredValueFromEnvAndValidate(dataDirKey, "required,dirpath")
	if err != nil {
		return Config{}, err
	}
	slog.Info("config", dataDirKey, dataDir)

	const logDirKey = "LOG_DIR"
	logDir := GetEnvWithDefaultValue(logDirKey, "required,dirpath", filepath.Join(dataDir, "logs"))
	slog.Info("config", logDirKey, logDir)

	const serverDirKey = "SERVER_DIR"
	serverDir := GetEnvWithDefaultValue(serverDirKey, "required,dirpath", filepath.Join(dataDir, "server"))
	slog.Info("config", serverDirKey, serverDir)

	const steamcmdDirKey = "STEAMCMD_DIR"
	steamcmdDir := GetEnvWithDefaultValue(steamcmdDirKey, "required,dirpath", filepath.Join(dataDir, "steamcmd"))
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
