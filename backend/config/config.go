package config

import (
	"cs-server-manager/gvalidator"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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

var errEnvNotFound = errors.New("environment variable not found or empty")

func getEnv(key string, validationString string) (string, error) {
	v, ok := os.LookupEnv(key)
	v = strings.TrimSpace(v)

	if !ok || v == "" {
		return "", errEnvNotFound
	}

	if err := gvalidator.Instance().Var(v, validationString); err != nil {
		return "", fmt.Errorf("environment variable with the key '%v' and the value '%v' failed to validate with validation string '%v': %w", key, v, validationString, err)
	}

	return v, nil
}

func getEnvWithDefaultValueIfEmptyAndIsDefaultIndicator(key string, validationString string, defaultValue string) (value string, isDefaultValue bool, err error) {
	v, err := getEnv(key, validationString)
	if err != nil {
		if errors.Is(err, errEnvNotFound) {
			return defaultValue, true, nil
		}
		return "", false, fmt.Errorf("validation of environment variable '%v' failed: %w", key, err)
	}

	return v, false, nil
}

func getEnvWithDefaultValueIfEmpty(key string, validationString string, defaultValue string) (string, error) {
	v, _, err := getEnvWithDefaultValueIfEmptyAndIsDefaultIndicator(key, validationString, defaultValue)
	if err != nil {
		return "", err
	}
	return v, nil
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

	// IP
	const ipKey = "IP"
	publicIp, err := GetPublicIp()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get public ip: %w", err)
	}
	ip, ipIsDefaultValue, err := getEnvWithDefaultValueIfEmptyAndIsDefaultIndicator(ipKey, "ip4_addr", publicIp)
	if err != nil {
		return Config{}, err
	}

	// HTTP_PORT
	const httpPortKey = "HTTP_PORT"
	httpPort, err := getEnvWithDefaultValueIfEmpty(httpPortKey, "port", "8080")
	if err != nil {
		return Config{}, err
	}

	// CS_PORT
	const csPortKey = "CS_PORT"
	csPort, err := getEnvWithDefaultValueIfEmpty(csPortKey, "port", "27015")
	if err != nil {
		return Config{}, err
	}

	// DATA_DIR
	workingDir, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get working directory: %w", err)
	}
	const dataDirKey = "DATA_DIR"
	dataDir, err := getEnvWithDefaultValueIfEmpty(dataDirKey, "dirpath", workingDir)
	if err != nil {
		return Config{}, err
	}

	dataDirAbs, err := filepath.Abs(dataDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to convert '%v' with value '%v' to absolute path: %w", dataDirKey, dataDir, err)
	}
	dataDir = dataDirAbs

	// LOG_DIR
	const logDirKey = "LOG_DIR"
	logDir, err := getEnvWithDefaultValueIfEmpty(logDirKey, "dirpath", filepath.Join(dataDir, "logs"))
	if err != nil {
		return Config{}, err
	}

	logDirAbs, err := filepath.Abs(logDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to convert '%v' with value '%v' to absolute path: %w", logDirKey, logDir, err)
	}
	logDir = logDirAbs

	// SERVER_DIR
	const serverDirKey = "SERVER_DIR"
	serverDir, err := getEnvWithDefaultValueIfEmpty(serverDirKey, "dirpath", filepath.Join(dataDir, "server"))
	if err != nil {
		return Config{}, err
	}

	serverDirAbs, err := filepath.Abs(serverDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to convert '%v' with value '%v' to absolute path: %w", serverDirKey, serverDir, err)
	}
	serverDir = serverDirAbs

	// STEAMCMD_DIR
	const steamcmdDirKey = "STEAMCMD_DIR"
	steamcmdDir, err := getEnvWithDefaultValueIfEmpty(steamcmdDirKey, "dirpath", filepath.Join(dataDir, "steamcmd"))
	if err != nil {
		return Config{}, err
	}

	steamcmdDirAbs, err := filepath.Abs(steamcmdDir)
	if err != nil {
		return Config{}, fmt.Errorf("failed to convert '%v' with value '%v' to absolute path: %w", steamcmdDirKey, steamcmdDir, err)
	}
	steamcmdDir = steamcmdDirAbs

	// ENABLE_WEB_UI
	const enableWebUiKey = "ENABLE_WEB_UI"
	enableWebUiStr, err := getEnvWithDefaultValueIfEmpty(enableWebUiKey, "boolean", "true")
	if err != nil {
		return Config{}, err
	}

	enableWebUi, err := strconv.ParseBool(enableWebUiStr)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse environment variable '%v' with value '%v' to bool: %w", enableWebUiKey, enableWebUiStr, err)
	}

	// ENABLE_SWAGGER
	const enableSwaggerKey = "ENABLE_SWAGGER"
	enableSwaggerStr, err := getEnvWithDefaultValueIfEmpty(enableSwaggerKey, "boolean", "true")
	if err != nil {
		return Config{}, err
	}

	enableSwagger, err := strconv.ParseBool(enableSwaggerStr)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse environment variable '%v' with value '%v' to bool: %w", enableSwaggerKey, enableSwaggerStr, err)
	}

	//
	cfg := Config{
		httpPort,
		csPort,
		dataDir,
		logDir,
		serverDir,
		steamcmdDir,
		enableWebUi,
		enableSwagger,
		ip,
		!ipIsDefaultValue,
	}

	// Print
	v := reflect.ValueOf(cfg)
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).IsExported() {
			continue
		}

		slog.Info("config", v.Type().Field(i).Name, v.Field(i).Interface())
	}

	return cfg, nil
}
