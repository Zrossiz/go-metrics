package memstorage

import (
	"testing"

	"github.com/Zrossiz/go-metrics/server/internal/storage"
)

func TestSetCounter(t *testing.T) {
	Metrics = append(Metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  "TestCounter",
		Value: int64(0),
	})

	SetCounter("TestCounter", 5)
	if val, ok := Metrics[0].Value.(int64); !ok || val != 5 {
		t.Fatalf("expected value 5, got %v", Metrics[0].Value)
	}

	// Проверяем SetCounter с вторым значением
	SetCounter("TestCounter", 5)
	if val, ok := Metrics[0].Value.(int64); !ok || val != 10 {
		t.Fatalf("expected value 10, got %v", Metrics[0].Value)
	}
}
