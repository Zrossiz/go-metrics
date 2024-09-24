package collector

import (
	"runtime"
	"testing"
)

func TestGetMetrics(t *testing.T) {
	var counter int64 = 1

	metrics := GetMetrics(&counter)

	for _, metric := range metrics {
		if metric.Name == "PollCount" {
			if val, ok := metric.Value.(int64); ok {
				if val != 2 {
					t.Errorf("expected poll count value 2, got %v", val)
				}
			} else {
				t.Errorf("expected metric.Value to be int, got %T", metric.Value)
			}
			break
		}
	}

	cpuCoresWithMetricsCount := runtime.NumCPU() + 31

	if len(metrics) != cpuCoresWithMetricsCount {
		t.Errorf("expected %v metrics, got %v", cpuCoresWithMetricsCount, len(metrics))
	}
}
