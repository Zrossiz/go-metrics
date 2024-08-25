package memstorage

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
)

func TestSetCounter(t *testing.T) {
	store := NewMemStorage()

	store.SetCounter("TestCounter", 5)
	if val, ok := store.Metrics[0].Value.(int64); !ok || val != 5 {
		t.Fatalf("expected value 5, got %v", store.Metrics[0].Value)
	}

	store.SetCounter("TestCounter", 5)
	if val, ok := store.Metrics[0].Value.(int64); !ok || val != 10 {
		t.Fatalf("expected value 10, got %v", store.Metrics[0].Value)
	}
}

func TestSetGauge(t *testing.T) {
	store := NewMemStorage()

	store.SetGauge("TestGauge", 0.5)
	if store.Metrics[0].Value != 0.5 {
		t.Fatalf("expected value 0.5, got %v", store.Metrics[0].Value)
	}

	store.SetGauge("TestGauge", 0.6)
	if store.Metrics[0].Value != 0.6 {
		t.Fatalf("expected value 0.6, got %v", store.Metrics[0].Value)
	}
}

func TestGetMetric(t *testing.T) {
	store := NewMemStorage()

	emptyMetric := store.GetMetric("TestEmpty")

	if emptyMetric != nil {
		t.Fatalf("expected empty value, got %v", emptyMetric.Name)
	}

	store.Metrics = append(store.Metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  "TestCounter",
		Value: 1,
	})

	metricCounter := store.GetMetric("TestCounter")
	if metricCounter.Value != 1 {
		t.Fatalf("expected value 1, got %v", metricCounter.Value)
	}

	store.Metrics = append(store.Metrics, storage.Metric{
		Type:  storage.GaugeType,
		Name:  "TestGauge",
		Value: 0.1,
	})

	metricGauge := store.GetMetric("TestGauge")
	if metricGauge.Value != 0.1 {
		t.Fatalf("expected value 0.1, got %v", metricGauge.Value)
	}
}
