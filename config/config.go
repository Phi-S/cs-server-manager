package config

import (
	"fmt"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort string
	CsPort   string
	DataDir  string
}

func GetConfig() (Config, error) {
	const envFile = ".env"
	if err := godotenv.Load(envFile); err != nil {
		return Config{}, fmt.Errorf("failed to load %q file", envFile)
	}

	const httpPortKey = "HTTP_PORT"
	httpPort, ok := os.LookupEnv(httpPortKey)
	if !ok || httpPort == "" {
		return Config{}, fmt.Errorf("failed to get %q from environment", httpPortKey)
	}
	if !govalidator.IsPort(httpPort) {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid PORT", httpPortKey, httpPort)
	}

	const csPortKey = "CS_PORT"
	csPort, ok := os.LookupEnv(csPortKey)
	if !ok || httpPort == "" {
		return Config{}, fmt.Errorf("failed to get %q from environment", csPortKey)
	}
	if !govalidator.IsPort(csPort) {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid PORT", csPortKey, csPort)
	}

	const dataDirKey = "DATA_DIR"
	dataDir, ok := os.LookupEnv(dataDirKey)
	if !ok || httpPort == "" {
		return Config{}, fmt.Errorf("failed to get %q from environment", dataDirKey)
	}
	if ok, _ := govalidator.IsFilePath(dataDir); !ok {
		return Config{}, fmt.Errorf("%q with the value %q is not a valid filepath", dataDirKey, dataDir)
	}

	return Config{
		httpPort,
		csPort,
		dataDir,
	}, nil
}
