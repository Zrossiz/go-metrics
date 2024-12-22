package app

import (
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/config"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/http/send"
	"github.com/Zrossiz/go-metrics/internal/agent/security"
	"github.com/Zrossiz/go-metrics/internal/agent/services/collector"
	"go.uber.org/zap"
)

const WorkerCount = 2

func StartAgent() {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatal("get config error", zap.Error(err))
	}

	fmt.Println("port: ", cfg.RunAddr)

	publicCryptoKey, err := security.GetPublicKey(cfg.PublicKeyPath)
	if err != nil {
		fmt.Println(err)
		zap.S().Fatal("get crypto key error", zap.Error(err))
	}
	cfg.PublicCryptoKey = publicCryptoKey

	pollIntervalDuration, err := time.ParseDuration(cfg.PollInterval)
	if err != nil {
		fmt.Println(err)
	}

	reportIntervalDuration, err := time.ParseDuration(cfg.ReportInterval)
	if err != nil {
		fmt.Println(err)
	}

	tickerPoll := time.NewTicker(pollIntervalDuration)
	tickerReport := time.NewTicker(reportIntervalDuration)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()

	metricsChan := make(chan []types.Metric, 10)
	sendChan := make(chan []types.Metric, 10)

	rateLimiter := make(chan struct{}, cfg.RateLimiter)

	for i := 0; i < WorkerCount; i++ {
		go collectorWorker(metricsChan)
	}

	for i := 0; i < WorkerCount; i++ {
		go senderWorker(sendChan, rateLimiter, cfg)
	}

	var counter int64

	for range tickerPoll.C {
		metrics := collector.GetMetrics(&counter)
		metricsChan <- metrics

		if len(tickerReport.C) > 0 {
			metrics := <-metricsChan
			sendChan <- metrics
		}
	}
}

func collectorWorker(metricsChan chan []types.Metric) {
	for metrics := range metricsChan {
		metricsChan <- metrics
	}
}

func senderWorker(sendChan chan []types.Metric, rateLimiter chan struct{}, cfg *config.Config) {
	for metrics := range sendChan {
		rateLimiter <- struct{}{}
		go func(metrics []types.Metric) {
			defer func() {
				<-rateLimiter
			}()
			send.Metrics(metrics, cfg)
		}(metrics)
	}
}
