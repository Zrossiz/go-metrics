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
	counter := 0.0

	for {
		select {
		case <-tickerPoll.C:
			metrics = collector.CollectMetrics()
			counter += 1
			fmt.Println("tick")
		case <-tickerReport.C:
			metrics = append(metrics, collector.Metric{
				Type:  collector.Counter,
				Name:  "PollCount",
				Value: counter,
			})
			handlers.SendMetrics(metrics)
			fmt.Println("report")
		}
	}
}

func main() {
	startMonitoring()
}
