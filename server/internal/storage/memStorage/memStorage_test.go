package memstorage

import (
	"testing"

	"github.com/Zrossiz/go-metrics/server/internal/storage"
)

func TestSetCounter(t *testing.T) {
	SetCounter("TestCounter", 5)
	if val, ok := Metrics[0].Value.(int64); !ok || val != 5 {
		t.Fatalf("expected value 5, got %v", Metrics[0].Value)
	}

	SetCounter("TestCounter", 5)
	if val, ok := Metrics[0].Value.(int64); !ok || val != 10 {
		t.Fatalf("expected value 10, got %v", Metrics[0].Value)
	}
}

func TestSetGauge(t *testing.T) {
	Metrics = []storage.Metric{}
	SetGauge("TestGauge", 0.5)
	if Metrics[0].Value != 0.5 {
		t.Fatalf("expected value 0.5, got %v", Metrics[0].Value)
	}

	SetGauge("TestGauge", 0.6)
	if Metrics[0].Value != 0.6 {
		t.Fatalf("expected value 0.6, got %v", Metrics[0].Value)
	}
}

func TestGetMetric(t *testing.T) {
	Metrics = []storage.Metric{}

	emptyMetric := GetMetric("TestEmpty")
	if emptyMetric.Name != "" {
		t.Fatalf("expected empty value, got %v", emptyMetric.Name)
	}

	Metrics = append(Metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  "TestCounter",
		Value: 1,
	})

	metricCounter := GetMetric("TestCounter")
	if metricCounter.Value != 1 {
		t.Fatalf("expected value 1, got %v", metricCounter.Value)
	}

	Metrics = append(Metrics, storage.Metric{
		Type:  storage.GaugeType,
		Name:  "TestGauge",
		Value: 0.1,
	})

	metricGauge := GetMetric("TestGauge")
	if metricGauge.Value != 0.1 {
		t.Fatalf("expected value 0.1, got %v", metricGauge.Value)
	}
}
