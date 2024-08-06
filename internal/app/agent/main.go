package agent

import (
	"fmt"
	"time"

	agentHandlers "github.com/Zrossiz/go-metrics/internal/http/agent"
	agentServices "github.com/Zrossiz/go-metrics/internal/services/agent"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func StartAgent() {
	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	defer tickerPoll.Stop()
	defer tickerReport.Stop()
	var metrics []agentServices.Metric
	counter := 0

	for {
		select {
		case <-tickerPoll.C:
			metrics = agentServices.CollectMetrics()
			counter += 1
			fmt.Println("tick")
		case <-tickerReport.C:
			metrics = append(metrics, agentServices.Metric{
				Type:  agentServices.Counter,
				Name:  "PollCount",
				Value: counter,
			})
			agentHandlers.SendMetrics(metrics)
			fmt.Println("report")
		}
	}
}
