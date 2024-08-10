package app

import (
	"fmt"
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
			send.SendMetrics(metrics, config.RunAddr)
			fmt.Println("report")
		}
	}
}
