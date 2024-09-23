package app

import (
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/config"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/http/send"
	"github.com/Zrossiz/go-metrics/internal/agent/services/collector"
)

const WorkerCount = 2

func StartAgent() {
	config.FlagParse()

	tickerPoll := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	metricsChan := make(chan []types.Metric, 10)
	sendChan := make(chan []types.Metric, 10)

	rateLimiter := make(chan struct{}, config.RateLimiter)

	for i := 0; i < WorkerCount; i++ {
		go collectorWorker(metricsChan)
	}

	for i := 0; i < WorkerCount; i++ {
		go senderWorker(sendChan, rateLimiter)
	}

	var counter int64 = 0

	for {
		select {
		case <-tickerPoll.C:
			metrics := collector.GetMetrics(&counter)
			metricsChan <- metrics

		case <-tickerReport.C:
			select {
			case metrics := <-metricsChan:
				sendChan <- metrics
			default:
				continue
			}
		}
	}
}

func collectorWorker(metricsChan chan []types.Metric) {
	for {
		select {
		case metrics := <-metricsChan:
			metricsChan <- metrics
		}
	}
}

func senderWorker(sendChan chan []types.Metric, rateLimiter chan struct{}) {
	for metrics := range sendChan {
		rateLimiter <- struct{}{}
		go func(metrics []types.Metric) {
			defer func() {
				<-rateLimiter
			}()

			send.Metrics(metrics, config.RunAddr)
		}(metrics)
	}
}
