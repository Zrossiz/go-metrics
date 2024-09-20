package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var RunAddr string
var PollInterval int64
var ReportInterval int64
var ApiKey string

func FlagParse() {
	_ = godotenv.Load()

	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&PollInterval, "p", 2, "interval for get metrics")
	flag.Int64Var(&ReportInterval, "r", 10, "interval for send metrics")

	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if val, err := strconv.ParseInt(envPollInterval, 2, 64); err == nil {
			PollInterval = val
		}
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.ParseInt(envReportInterval, 6, 64); err == nil {
			ReportInterval = val
		}
	}

	if envApiKey := os.Getenv("KEY"); envApiKey != "" {
		ApiKey = envApiKey
	} else {
		flag.StringVar(&ApiKey, "k", "", "api key for hash")
	}
}
