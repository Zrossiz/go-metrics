// Config parser package
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
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress   string `json:"address"`
	StoreInterval   string `json:"store_interval"`
	FileStoragePath string
	Restore         bool   `json:"restore"`
	DBDSN           string `json:"database_dsn"`
	LogLevel        string `json:"log_level"`
	HashKey         string `json:"hash_key"`
	PrivateKeyPath  string `json:"crypto_key"`
	TrustedSubnet   string `json:"trusted_subnet"`
	JSONConfigPath  string
	PrivateKey      *rsa.PrivateKey
}

var AppConfig Config

// Parsing config
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
		fmt.Print("Start parsing config from JSON file...\n")
		cfgJSON, err := loadJSONConfig(cfg.JSONConfigPath)
		if err != nil {
			fmt.Printf("Error parsing config from JSON: %v\n", err)
		}

		cfg = cfgJSON
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "address and port to run server")
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		cfg.StoreInterval = envStoreInterval
	} else {
		flag.StringVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "interval for save metrics")
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
		cfg.Restore = value
	} else {
		flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "get metrics from file")

	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "path to storage file")
	}

	if envDBConn := os.Getenv("DATABASE_DSN"); envDBConn != "" {
		cfg.DBDSN = envDBConn
	} else {
		flag.StringVar(&cfg.DBDSN, "d", cfg.DBDSN, "dsn for database")
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.StringVar(&cfg.PrivateKeyPath, "crypto-key", cfg.PrivateKeyPath, "private key for encrypt request")
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cfg.PrivateKeyPath = envCryptoKey
	}

	flag.StringVar(&cfg.HashKey, "k", cfg.HashKey, "key for hash")
	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.HashKey = envKey
	}

	flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "trusted subnet")
	if envTrustedSubNet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubNet != "" {
		cfg.TrustedSubnet = envTrustedSubNet
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
