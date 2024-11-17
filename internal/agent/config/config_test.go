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

func TestGetConfig_InvalidEnvValues(t *testing.T) {
	resetFlags()

	os.Setenv("POLL_INTERVAL", "invalid")
	os.Setenv("REPORT_INTERVAL", "notanumber")
	os.Setenv("RATE_LIMITER", "badlimit")

	cfg, err := GetConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
