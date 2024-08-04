package config_test

import (
	"cs-server-manager/config"
	globalvalidator "cs-server-manager/global_validator"
	"os"
	"testing"
)

func init() {
	globalvalidator.Init()
}

func TestGetRequiredValueFromEnvAndValidate(t *testing.T) {
	testEnvKey := "test-env-key"
	testEnvValue := "65535"

	os.Setenv(testEnvKey, testEnvValue)
	defer os.Unsetenv(testEnvKey)

	value, err := config.GetRequiredValueFromEnvAndValidate(testEnvKey, "required,port")
	if err != nil {
		t.Error(err)
	} else if value != testEnvValue {
		t.Error("returned value dose not match expected value")
	}
}
