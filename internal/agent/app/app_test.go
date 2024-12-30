package app

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func TestCollectorWorker(t *testing.T) {
	metricsChan := make(chan []types.Metric, 1)
	go collectorWorker(metricsChan)

	inputMetrics := []types.Metric{{Name: "test_metric", Value: 1}}
	metricsChan <- inputMetrics

	result := <-metricsChan
	if len(result) != len(inputMetrics) || result[0].Name != "test_metric" {
		t.Errorf("expected %v, got %v", inputMetrics, result)
	}
}
