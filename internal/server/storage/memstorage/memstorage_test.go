package memstorage

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	storage := New()

	assert.NotNil(t, storage)
	assert.Empty(t, storage.data)
}

func TestMemStorageSetGauge(t *testing.T) {
	storage := New()

	value := 123.45
	metric := dto.PostMetricDto{
		ID:    "TestGauge",
		MType: models.GaugeType,
		Value: &value,
	}

	storage.SetGauge(metric)

	assert.Len(t, storage.data, 1)

	addedMetric := storage.data[0]
	assert.Equal(t, metric.ID, addedMetric.Name)
	assert.Equal(t, metric.MType, addedMetric.Type)
	assert.Equal(t, *metric.Value, *addedMetric.Value)
}

func TestMemStorageSetGauge_UpdateExisting(t *testing.T) {
	storage := New()

	value1 := 123.45
	metric1 := dto.PostMetricDto{
		ID:    "TestGauge",
		MType: models.GaugeType,
		Value: &value1,
	}

	storage.SetGauge(metric1)

	value2 := 678.90
	metric2 := dto.PostMetricDto{
		ID:    "TestGauge",
		MType: models.GaugeType,
		Value: &value2,
	}
	storage.SetGauge(metric2)

	assert.Len(t, storage.data, 1)
	addedMetric := storage.data[0]
	assert.Equal(t, value2, *addedMetric.Value)
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

	assert.Len(t, storage.data, 1)

	addedMetric := storage.data[0]
	assert.Equal(t, metric.ID, addedMetric.Name)
	assert.Equal(t, metric.MType, addedMetric.Type)
	assert.Equal(t, *metric.Delta, *addedMetric.Delta)
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

	assert.NotNil(t, addedMetric)
	assert.Equal(t, name, addedMetric.Name)
	assert.Equal(t, *metric.Delta, *addedMetric.Delta)
}

func TestGetMetric_Empty(t *testing.T) {
	storage := New()

	metric, err := storage.Get("not_found")
	assert.NoError(t, err)
	assert.Nil(t, metric)
}

func TestGetAllMetrics(t *testing.T) {
	storage := New()

	var metricValue1 int64 = 123
	metricValue2 := 1.2

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
	assert.Len(t, addedMetrics, 2)
}

func TestGetAllMetrics_Empty(t *testing.T) {
	storage := New()

	metrics, _ := storage.GetAll()
	assert.Empty(t, metrics)
}

func TestSetBatch(t *testing.T) {
	storage := New()

	var metricValue1 int64 = 123
	metricValue2 := 1.2

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
	assert.Len(t, addedMetrics, 2)
}

func TestSetBatch_UpdateExisting(t *testing.T) {
	storage := New()

	var metricValue1 int64 = 123
	metricValue2 := 1.2
	metricValue3 := 999.99

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

	metrics[1].Value = &metricValue3
	_ = storage.SetBatch([]dto.PostMetricDto{metrics[1]})

	addedMetrics, _ := storage.GetAll()
	assert.Len(t, addedMetrics, 2)
	assert.Equal(t, metricValue3, *addedMetrics[1].Value)
}
