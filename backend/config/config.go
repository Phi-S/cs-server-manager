package config

import (
	"cs-server-manager/gvalidator"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort      string
	CsPort        string
	DataDir       string
	LogDir        string
	ServerDir     string
	SteamcmdDir   string
	EnableWebUi   bool
	EnableSwagger bool
	Ip            string
	usedEnvIp     bool
}

func (c Config) UsedIpFromEnv() bool {
	return c.usedEnvIp
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

	if err := gvalidator.Instance().Var(value, validationString); err != nil {
		slog.Warn("validation failed. Returning default value", "validation_string", validationString, "env_key", key, "default_value", defaultValue)
		return defaultValue
	}

	return value
}

func GetEnvWithDefaultValue2(key string, validationString string, defaultValue string) (value string, isDefaultValue bool) {
	value, ok := os.LookupEnv(key)
	value = strings.TrimSpace(value)

	if !ok || value == "" {
		return defaultValue, true
	}

	if err := gvalidator.Instance().Var(value, validationString); err != nil {
		slog.Warn("validation failed. Returning default value", "validation_string", validationString, "env_key", key, "default_value", defaultValue)
		return defaultValue, true
	}

	return value, false
}

func GetPublicIp() (string, error) {
	resp, err := http.Get("https://api.ipify.org/?format=text")
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body %w", err)
	}

	return string(body), nil
}

func GetConfig() (Config, error) {
	const envFile = ".env"
	if err := godotenv.Load(envFile); err != nil {
		slog.Info("no .env file present at", "path", envFile)
	}

	const ipKey = "IP"
	publicIp, err := GetPublicIp()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get public ip: %w", err)
	}
	ip, isDefaultValue := GetEnvWithDefaultValue2(ipKey, "ip4_addr", publicIp)
	slog.Info("config", ipKey, ip)

	const httpPortKey = "HTTP_PORT"
	httpPort := GetEnvWithDefaultValue(httpPortKey, "port", "8080")
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

	const enableWebUiKey = "ENABLE_WEB_UI"
	enableWebUi := GetEnvWithDefaultValueBool(enableWebUiKey, "boolean", true)
	slog.Info("config", enableWebUiKey, enableWebUi)

	const enableSwaggerKey = "ENABLE_SWAGGER"
	enableSwagger := GetEnvWithDefaultValueBool(enableSwaggerKey, "boolean", true)
	slog.Info("config", enableSwaggerKey, enableSwagger)

	return Config{
		httpPort,
		csPort,
		dataDir,
		logDir,
		serverDir,
		steamcmdDir,
		enableWebUi,
		enableSwagger,
		ip,
		!isDefaultValue,
	}, nil
}
