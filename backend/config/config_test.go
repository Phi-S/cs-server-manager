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

const IP_KEY = "IP"

func Test_ipSetByEnvironmentVariable_false(t *testing.T) {
	ip := "127.127.127.127"
	err := os.Setenv(IP_KEY, ip)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(IP_KEY)

	config, err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	if !config.ipSetByEnvironmentVariable {
		t.Fatal("ipSetByEnvironmentVariable is false but should be true")
	}

	if config.Ip != ip {
		t.Fatal("ip in config dose not match the one from GetConfig")
	}
}

func Test_ipSetByEnvironmentVariable_true(t *testing.T) {
	os.Unsetenv(IP_KEY)

	config, err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	if config.ipSetByEnvironmentVariable {
		t.Fatal("ipSetByEnvironmentVariable is true but should be false")
	}
}
