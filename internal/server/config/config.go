package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress   string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	DATABASE_DSN    string
	LogLevel        string
}

var AppConfig Config

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerAddress = envRunAddr
	} else {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return nil, err
		}
		cfg.StoreInterval = value
	} else {
		flag.IntVar(&cfg.StoreInterval, "i", 5, "interval for save metrics")
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			return nil, err
		}
		cfg.Restore = value
	} else {
		flag.BoolVar(&cfg.Restore, "r", false, "get metrics from file")
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(&cfg.FileStoragePath, "f", "", "path to storage file")
	}

	if envDBConn := os.Getenv("DB_DSN"); envDBConn != "" {
		cfg.DATABASE_DSN = envDBConn
	} else {
		flag.StringVar(&cfg.DATABASE_DSN, "d", "", "dsn for database")
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	} else {
		cfg.LogLevel = zapcore.ErrorLevel.String()
	}

	flag.Parse()

	return cfg, nil
}
