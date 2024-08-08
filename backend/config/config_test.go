package config_test

import (
	"cs-server-manager/config"
	globalvalidator "cs-server-manager/gvalidator"
	"os"
	"testing"
)

func init() {
	globalvalidator.Init()
}

func TestGetEnvWithDefaultValue_port_OK(t *testing.T) {
	testEnvKey := "test-env-key"
	testEnvValue := "65535"

	err := os.Setenv(testEnvKey, testEnvValue)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(testEnvKey)

	value := config.GetEnvWithDefaultValue(testEnvKey, "port", "1234")
	if value != testEnvValue {
		t.Fatalf("unexpected value received %v. expected value: %v", value, testEnvValue)
	}
}
