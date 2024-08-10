package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var RunAddr string

func FlagParse() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка при загрузке файла .env:", err)
	} else {
		log.Println("Файл .env загружен успешно")
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		RunAddr = envRunAddr
	} else {
		flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()
}
