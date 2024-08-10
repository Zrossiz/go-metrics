package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var RunAddr string
var PollInterval int64
var ReportInterval int64

func FlagParse() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка при загрузке файла .env:", err)
	} else {
		log.Println("Файл .env загружен успешно")
	}

	if envRunAddr := os.Getenv("ADDRESS_SERVER"); envRunAddr != "" {
		RunAddr = envRunAddr
	} else {
		flag.StringVar(&RunAddr, "a", "localhost:8080", "address and port to run server")
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {

		if val, err := strconv.ParseInt(envPollInterval, 10, 64); err == nil {
			PollInterval = val
		}
	} else {
		flag.Int64Var(&PollInterval, "p", 2, "interval for get metrics")
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.ParseInt(envReportInterval, 10, 64); err == nil {
			ReportInterval = val
		}
	} else {
		flag.Int64Var(&ReportInterval, "r", 10, "interval for send metrics")
	}
}
