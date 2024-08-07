package app

import (
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/http"
	"github.com/Zrossiz/go-metrics/internal/agent/services/collector"
)

func StartAgent() {
	tickerPoll := time.NewTicker(constants.PollInterval)
	tickerReport := time.NewTicker(constants.ReportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()
	var metrics []types.Metric
	counter := 0

	for {
		select {
		case <-tickerPoll.C:
			metrics = collector.CollectMetrics()
			counter += 1
			fmt.Println("tick")
		case <-tickerReport.C:
			metrics = append(metrics, types.Metric{
				Type:  constants.Counter,
				Name:  "PollCount",
				Value: counter,
			})
			http.SendMetrics(metrics)
			fmt.Println("report")
		}
	}
}
