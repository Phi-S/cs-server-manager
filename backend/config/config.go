package config

import (
	globalvalidator "cs-server-manager/gvalidator"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort       string
	CsPort         string
	DataDir        string
	LogDir         string
	ServerDir      string
	SteamcmdDir    string
	DisableWebsite bool
}

func GetEnvWithDefaultValueBool(key string, validationString string, defaultValue bool) bool {
	defaultValueStr := strconv.FormatBool(defaultValue)
	value := GetEnvWithDefaultValue(key, validationString, defaultValueStr)

	if value == defaultValueStr {
		return defaultValue
	}

	valueBool, err := strconv.ParseBool(value)
	if err != nil {
		slog.Warn("failed to parse value from env as bool", "value", value)
		return defaultValue
	}

	return valueBool
}

func GetEnvWithDefaultValue(key string, validationString string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return defaultValue
	}

	if err := globalvalidator.Instance.Var(value, validationString); err != nil {
		slog.Warn("validation failed. Returning default value", "validation_string", validationString, "env_key", key, "default_value", defaultValue)
		return defaultValue
	}

	return value
}

func GetConfig() (Config, error) {

	const envFile = ".env"
	if err := godotenv.Load(envFile); err != nil {
		slog.Info("no .env file present at", "path", envFile)
	}

	const httpPortKey = "HTTP_PORT"
	httpPort := GetEnvWithDefaultValue(httpPortKey, "port", "80")
	slog.Info("config", httpPortKey, httpPort)

	const csPortKey = "CS_PORT"
	csPort := GetEnvWithDefaultValue(csPortKey, "port", "27015")
	slog.Info("config", csPortKey, csPort)

	const dataDirKey = "DATA_DIR"
	dataDir := GetEnvWithDefaultValue(dataDirKey, "dirpath", "/data")
	slog.Info("config", dataDirKey, dataDir)

	const logDirKey = "LOG_DIR"
	logDir := GetEnvWithDefaultValue(logDirKey, "dirpath", filepath.Join(dataDir, "logs"))
	slog.Info("config", logDirKey, logDir)

	const serverDirKey = "SERVER_DIR"
	serverDir := GetEnvWithDefaultValue(serverDirKey, "dirpath", filepath.Join(dataDir, "server"))
	slog.Info("config", serverDirKey, serverDir)

	const steamcmdDirKey = "STEAMCMD_DIR"
	steamcmdDir := GetEnvWithDefaultValue(steamcmdDirKey, "dirpath", filepath.Join(dataDir, "steamcmd"))
	slog.Info("config", steamcmdDirKey, steamcmdDir)

	const disableWebsiteKey = "DISABLE_WEBSITE"
	disableWebsite := GetEnvWithDefaultValueBool(disableWebsiteKey, "boolean", false)
	slog.Info("config", disableWebsiteKey, disableWebsite)

	return Config{
		httpPort,
		csPort,
		dataDir,
		logDir,
		serverDir,
		steamcmdDir,
		disableWebsite,
	}, nil
}
