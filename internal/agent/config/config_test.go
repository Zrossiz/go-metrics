package config

import (
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

	os.Setenv("POLL_INTERVAL", "10")
	os.Setenv("REPORT_INTERVAL", "10")
	os.Setenv("RATE_LIMITER", "10")
	os.Setenv("KEY", "hash1")
	os.Setenv("ADDRESS", "addr")

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestGetConfig_ValidJSONConfig(t *testing.T) {
	resetFlags()

	JSONConfigFilePath, err := createFileWithJSONConfig()
	if err != nil {
		t.Errorf("create json config file: %v", err)
	}

	//init config

	err = deleteFileWithJSONConfig(JSONConfigFilePath)
	if err != nil {
		t.Errorf("delete json config file: %v", err)
	}

}

func createFileWithJSONConfig() (string, error) {
	return "", nil
}

func deleteFileWithJSONConfig(string) error {
	return nil
}

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}
