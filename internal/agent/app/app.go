package app

import (
	"flag"
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/http"
	"github.com/Zrossiz/go-metrics/internal/agent/services/collector"
)

var ServerAddress = flag.String("a", "localhost:8080", "Server address")
var PollInterval = flag.Int("p", 2, "Poll interval")
var ReportInterval = flag.Int("r", 10, "Report interval")

func StartAgent() {
	flag.Parse()

	tickerPoll := time.NewTicker(time.Duration(*PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(*ReportInterval) * time.Second)
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
			addr := fmt.Sprintf("%v", *ServerAddress)
			http.SendMetrics(metrics, addr)
			fmt.Println("report")
		}
	}
}
