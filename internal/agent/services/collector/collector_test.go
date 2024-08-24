package collector

import (
	"testing"
)

func TestGettMetrics(t *testing.T) {
	metrics := GetMetrics()

	if len(metrics) == 0 {
		t.Fatal("Expected non empty slice")
	}

	expectedMetrics := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
		"Lookups", "MCacheInuse", "MCacheSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys",
		"TotalAlloc", "RandomValue",
	}

	for i := 0; i < len(expectedMetrics); i++ {
		found := false

		for _, metric := range metrics {
			if expectedMetrics[i] == metric.Name {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("expected %s, got undefined", expectedMetrics[i])
		}
	}
}
