package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RunAddr        string
	PollInterval   int64
	ReportInterval int64
	Key            string
	RateLimiter    int64
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "interval for get metrics")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "interval for send metrics")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.RunAddr = envRunAddr
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		val, err := strconv.ParseInt(envPollInterval, 2, 64)

		if err == nil {
			cfg.PollInterval = val
		} else {
			return nil, err
		}
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		val, err := strconv.ParseInt(envReportInterval, 6, 64)
		if err == nil {
			cfg.ReportInterval = val
		} else {
			return nil, err
		}
	}

	flag.StringVar(&cfg.Key, "k", "", "key for hash")
	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.Key = envKey
	}

	flag.Int64Var(&cfg.RateLimiter, "l", 1000, "rate limiter")
	if envRateLimiter := os.Getenv("RATE_LIMITER"); envRateLimiter != "" {
		value, err := strconv.Atoi(envRateLimiter)
		if err == nil {
			cfg.RateLimiter = int64(value)
		} else {
			return nil, err
		}
	}

	flag.Parse()

	return cfg, nil
}
