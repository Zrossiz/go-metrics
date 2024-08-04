package main

import (
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/http/handlers"
	"github.com/Zrossiz/go-metrics/internal/agent/lib/collector"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func startMonitoring() {
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()
	var metrics []collector.Metric

	for {
		select {
		case <-tickerPoll.C:
			metrics = collector.CollectMetrics()
			fmt.Println("tick")
		case <-tickerReport.C:
			handlers.SendMetrics(metrics)
			fmt.Println("report")
		}
	}
}

func main() {
	startMonitoring()
}
