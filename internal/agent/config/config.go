package config

import (
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RunAddr         string `json:"address"`
	PollInterval    string `json:"poll_interval"`
	ReportInterval  string `json:"report_interval"`
	HashKey         string
	RateLimiter     int64
	PublicKeyPath   string `json:"crypto_key"`
	JSONConfigPath  string
	PublicCryptoKey *rsa.PublicKey
}

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if envJSONConfigPath := os.Getenv("CONFIG"); envJSONConfigPath != "" {
		cfg.JSONConfigPath = envJSONConfigPath
	} else {
		flag.StringVar(&cfg.JSONConfigPath, "config", "", "file path for json config")
	}
	flag.Parse()

	if cfg.JSONConfigPath != "" {
		cfgJSON, err := loadJSONConfig(cfg.JSONConfigPath)
		if err != nil {
			fmt.Printf("Eror parsing config from JSON: %v\n", err)
		}

		cfg = cfgJSON
	}

	flag.StringVar(&cfg.PollInterval, "p", "2s", "interval for get metrics")
	flag.StringVar(&cfg.ReportInterval, "r", "4s", "interval for send metrics")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.RunAddr = envRunAddr
	} else {
		flag.StringVar(&cfg.RunAddr, "a", cfg.RunAddr, "address and port to run server")
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		cfg.PollInterval = envPollInterval
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		cfg.ReportInterval = envReportInterval
	}

	flag.StringVar(&cfg.HashKey, "k", "", "key for hash")
	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.HashKey = envKey
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

	flag.StringVar(&cfg.PublicKeyPath, "crypto-key", "/Users/zrossiz/Desktop/GoProjects/praktikum/projects/go-metrics/crypto/public_key.pem", "public key for encrypt request")
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cfg.PublicKeyPath = envCryptoKey
	}

	flag.Parse()

	return cfg, nil
}

func loadJSONConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(fileContent, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
