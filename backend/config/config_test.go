package config

import (
	"os"
	"testing"
)

func Test_getEnvWithDefaultValueIfEmpty_OK(t *testing.T) {
	testEnvKey := "test-env-key"
	testEnvValue := "65535"

	err := os.Setenv(testEnvKey, testEnvValue)
	if err != nil {
		t.Fatal(err)
	}
	defer func(key string) {
		_ = os.Unsetenv(key)
	}(testEnvKey)

	value, err := getEnvWithDefaultValueIfEmpty(testEnvKey, "port", "1234")
	if err != nil {
		t.Fatal(err)
	}

	if value != testEnvValue {
		t.Fatalf("unexpected value received %v. expected value: %v", value, testEnvValue)
	}
}

func Test_getEnvWithDefaultValueIfEmpty_DefaultValue(t *testing.T) {
	testEnvKey := "test-env-key"
	defaultValue := "12345"

	defer func(key string) {
		_ = os.Unsetenv(key)
	}(testEnvKey)

	value, err := getEnvWithDefaultValueIfEmpty(testEnvKey, "port", defaultValue)
	if err != nil {
		t.Fatal(err)
	}

	if value != defaultValue {
		t.Fatalf("unexpected value received '%v'. expected value: '%v'", value, defaultValue)
	}
}
