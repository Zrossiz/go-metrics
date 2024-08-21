package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

var (
	RunAddr      string
	FlagLogLevel string
)

func FlagParse() {
	_ = godotenv.Load()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	} else {
		flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	} else {
		FlagLogLevel = zapcore.ErrorLevel.String()
	}
}
