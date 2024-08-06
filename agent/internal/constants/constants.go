package constants

import "time"

const (
	Counter        = "counter"
	Gauge          = "gauge"
	PollInterval   = 2 * time.Second
	ReportInterval = 10 * time.Second
)

const ServerAddress = "localhost:8080"
