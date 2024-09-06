package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

var (
	RunAddr         string
	StoreInterval   int
	Restore         bool
	FileStoragePath string
	FlagLogLevel    string
	DbConnString    string
)

func FlagParse() error {
	_ = godotenv.Load()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	} else {
		flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		value, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			return err
		}
		StoreInterval = value
	} else {
		flag.IntVar(&StoreInterval, "i", 20, "interval for save metrics")
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		value, err := strconv.ParseBool(envRestore)
		if err != nil {
			return err
		}
		Restore = value
	} else {
		flag.BoolVar(&Restore, "r", false, "get metrics from file")
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	} else {
		flag.StringVar(&FileStoragePath, "f", "storage/storage.txt", "path to storage file")
	}

	if envDbConn := os.Getenv("DB_DSN"); envDbConn != "" {
		DbConnString = envDbConn
	} else {
		flag.StringVar(&DbConnString, "d", "", "dsn for database")
	}

	flag.Parse()

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	} else {
		FlagLogLevel = zapcore.ErrorLevel.String()
	}

	return nil
}
