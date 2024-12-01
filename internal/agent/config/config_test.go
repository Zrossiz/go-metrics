package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestGetConfig_Success(t *testing.T) {
	resetFlags()
	os.Args = []string{
		"cmd",
		"-a", "192.168.1.1:8080",
		"-p", "4",
		"-r", "11",
		"-l", "1001",
		"-k", "hash1",
	}

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, 4, int(cfg.PollInterval))
	assert.Equal(t, 11, int(cfg.ReportInterval))
	assert.Equal(t, "hash1", cfg.Key)
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
