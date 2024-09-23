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
var Key string
var RateLimiter int64

func FlagParse() {
	_ = godotenv.Load()

	flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&PollInterval, "p", 2, "interval for get metrics")
	flag.Int64Var(&ReportInterval, "r", 10, "interval for send metrics")

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

	flag.StringVar(&Key, "k", "", "key for hash")
	if envKey := os.Getenv("KEY"); envKey != "" {
		Key = envKey
	}

	flag.Int64Var(&RateLimiter, "l", 1000, "rate limiter")
	if envRateLimiter := os.Getenv("RATE_LIMITER"); envRateLimiter != "" {
		value, err := strconv.Atoi(envRateLimiter)
		if err == nil {
			RateLimiter = int64(value)
		}
	}

	flag.Parse()
}
