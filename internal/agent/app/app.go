package app

import (
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/config"
	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/http/send"
	"github.com/Zrossiz/go-metrics/internal/agent/services/collector"
)

func StartAgent() {
	config.FlagParse()

	tickerPoll := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()
	var metrics []types.Metric
	var counter int64

	for {
		select {
		case <-tickerPoll.C:
			metrics = collector.GetMetrics()
			counter++
		case <-tickerReport.C:
			metrics = append(metrics, types.Metric{
				Type:  constants.Counter,
				Name:  "PollCount",
				Value: counter,
			})
			send.GzipMetrics(metrics, config.RunAddr)
		}
	}
}
