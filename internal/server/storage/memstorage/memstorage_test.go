package memstorage

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
)

func TestMemStorageSetGauge(t *testing.T) {
	storage := New()

	value := 123.45

	metric := dto.PostMetricDto{
		ID:    "TestGauge",
		MType: models.GaugeType,
		Value: &value,
	}

	storage.SetGauge(metric)

	if len(storage.data) == 0 {
		t.Errorf("expected 1 element, got 0")
	}

	addedMetric := storage.data[0]

	if addedMetric.Value != &value {
		t.Errorf("exptected %v, got %v", value, *addedMetric.Value)
	}

	if addedMetric.Type != models.GaugeType {
		t.Errorf("expected gauge type, got %v", addedMetric.Type)
	}
}

func TestMemStorageSetCounter(t *testing.T) {
	storage := New()

	var value int64 = 123

	metric := dto.PostMetricDto{
		ID:    "TestCounter",
		MType: models.CounterType,
		Delta: &value,
	}

	storage.SetCounter(metric)

	if len(storage.data) == 0 {
		t.Errorf("expected 1 element, got 0")
	}

	addedMetric := storage.data[0]

	if addedMetric.Delta != &value {
		t.Errorf("exptected %v, got %v", value, *addedMetric.Value)
	}

	if addedMetric.Type != models.CounterType {
		t.Errorf("expected gauge type, got %v", addedMetric.Type)
	}
}

func TestGetMetric(t *testing.T) {
	storage := New()

	var value int64 = 123
	name := "TestCounter"

	metric := dto.PostMetricDto{
		ID:    name,
		MType: models.CounterType,
		Delta: &value,
	}

	storage.SetCounter(metric)

	addedMetric, _ := storage.Get(name)
	if addedMetric == nil {
		t.Errorf("expected metric, got nil")
	}

	if addedMetric.Name != name {
		t.Errorf("expected name %v, got %v", name, addedMetric.ID)
	}

	if addedMetric.Delta != &value {
		t.Errorf("exptected %v, got %v", value, *addedMetric.Value)
	}
}

func TestGetAllMetrics(t *testing.T) {
	storage := New()

	var metricValue1 int64 = 123
	var metricValue2 float64 = 1.2

	metrics := []dto.PostMetricDto{
		{
			ID:    "TestCounter",
			MType: models.CounterType,
			Delta: &metricValue1,
		},
		{
			ID:    "TestGauge",
			MType: models.GaugeType,
			Value: &metricValue2,
		},
	}

	storage.SetCounter(metrics[0])
	storage.SetGauge(metrics[1])

	addedMetrics, _ := storage.GetAll()
	if len(*addedMetrics) != 2 {
		t.Errorf("expected 2 metrics, got %v", len(*addedMetrics))
	}
}

func TestSetBatch(t *testing.T) {
	storage := New()

	var metricValue1 int64 = 123
	var metricValue2 float64 = 1.2

	metrics := []dto.PostMetricDto{
		{
			ID:    "TestCounter",
			MType: models.CounterType,
			Delta: &metricValue1,
		},
		{
			ID:    "TestGauge",
			MType: models.GaugeType,
			Value: &metricValue2,
		},
	}

	_ = storage.SetBatch(metrics)

	addedMetrics, _ := storage.GetAll()
	if len(*addedMetrics) != 2 {
		t.Errorf("expected 2 metrics, got %v", len(*addedMetrics))
	}
}
