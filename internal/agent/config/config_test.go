package config

import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_Success(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-a", "192.168.1.1:8080",
		"-p", "4s",
		"-r", "11s",
		"-l", "1001",
		"-k", "hash1",
	}

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "4s", cfg.PollInterval)
	assert.Equal(t, "11s", cfg.ReportInterval)
	assert.Equal(t, "hash1", cfg.HashKey)
	assert.Equal(t, "192.168.1.1:8080", cfg.RunAddr)
	assert.Equal(t, 1001, int(cfg.RateLimiter))
}

func TestGetConfig_InvalidEnvValues(t *testing.T) {
	resetFlags()

	os.Setenv("POLL_INTERVAL", "invalid")
	os.Setenv("REPORT_INTERVAL", "notanumber")
	os.Setenv("RATE_LIMITER", "badlimit")

	cfg, err := GetConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestGetConfig_ValidEnvValues(t *testing.T) {
	resetFlags()

	os.Setenv("POLL_INTERVAL", "10s")
	os.Setenv("REPORT_INTERVAL", "10s")
	os.Setenv("RATE_LIMITER", "10")
	os.Setenv("KEY", "hash1")
	os.Setenv("ADDRESS", "addr")

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestGetConfig_ValidJSONConfig(t *testing.T) {
	resetFlags()
	os.Clearenv()

	JSONConfigFilePath, err := createFileWithJSONConfig()
	if err != nil {
		t.Errorf("create json config file: %v", err)
	}
	defer deleteFileWithJSONConfig(JSONConfigFilePath)

	os.Setenv("CONFIG", JSONConfigFilePath)

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "4s", cfg.PollInterval)
	// assert.Equal(t, "11s", cfg.ReportInterval)
	// assert.Equal(t, "hash1", cfg.HashKey)
	// assert.Equal(t, "192.168.1.1:8080", cfg.RunAddr)
	// assert.Equal(t, 1001, int(cfg.RateLimiter))
}

func createFileWithJSONConfig() (string, error) {
	tmpFile, err := os.CreateTemp("", "config_*.json")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	cfg := &Config{
		RunAddr:        "192.168.1.1:8080",
		PollInterval:   "4s",
		ReportInterval: "11s",
		HashKey:        "hash1",
		RateLimiter:    1001,
	}

	fileContent, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}

	if _, err := tmpFile.Write(fileContent); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func deleteFileWithJSONConfig(filePath string) error {
	return os.Remove(filePath)
}

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}
