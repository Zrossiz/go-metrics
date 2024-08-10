package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

var RunAddr string

func FlagParse() {
	_ = godotenv.Load()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	} else {
		flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()
}
