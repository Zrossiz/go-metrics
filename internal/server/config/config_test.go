package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestGetConfig_FlagOverrides(t *testing.T) {
	resetFlags()

	os.Clearenv()

	os.Args = []string{
		"cmd",
		"-a", "192.168.1.1:8080",
		"-i", "15",
		"-r",
		"-f", "/var/metrics.json",
		"-d", "postgres://user:pass@localhost/dbname",
		"-k", "overridekey",
	}

	cfg, err := GetConfig()
	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, "192.168.1.1:8080", cfg.ServerAddress)
	assert.Equal(t, 15, cfg.StoreInterval)
	assert.True(t, cfg.Restore)
	assert.Equal(t, "/var/metrics.json", cfg.FileStoragePath)
	assert.Equal(t, "postgres://user:pass@localhost/dbname", cfg.DBDSN)
	assert.Equal(t, "error", cfg.LogLevel)
	assert.Equal(t, "overridekey", cfg.HashKey)
}

func TestGetConfig_InvalidEnvValues(t *testing.T) {
	resetFlags()

	os.Setenv("STORE_INTERVAL", "invalid")
	os.Setenv("RESTORE", "notabool")

	cfg, err := GetConfig()
	assert.Error(t, err, "Error should occur with invalid env values")
	assert.Nil(t, cfg, "Config should be nil with invalid env values")
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
