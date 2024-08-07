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
