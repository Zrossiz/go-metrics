package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T) *os.File {
	tempFile, err := os.CreateTemp("", "filestorage_test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tempFile
}

func TestSetGauge(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	metric := dto.PostMetricDto{ID: "gauge_metric", MType: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()}

	err := storage.SetGauge(metric)
	assert.NoError(t, err)

	storedMetric, err := storage.Get("gauge_metric")
	assert.NoError(t, err)
	assert.NotNil(t, storedMetric)
	assert.Equal(t, metric.Value, storedMetric.Value)
}

func TestSetCounter(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	initialDelta := int64(10)
	metric := dto.PostMetricDto{ID: "counter_metric", MType: models.CounterType, Delta: &initialDelta}

	err := storage.SetCounter(metric)
	assert.NoError(t, err)

	storedMetric, err := storage.Get("counter_metric")
	assert.NoError(t, err)
	assert.NotNil(t, storedMetric)
	assert.Equal(t, initialDelta, *storedMetric.Delta)

	// Test incrementing the counter
	additionalDelta := int64(5)
	metric.Delta = &additionalDelta
	err = storage.SetCounter(metric)
	assert.NoError(t, err)

	storedMetric, err = storage.Get("counter_metric")
	assert.NoError(t, err)
	assert.Equal(t, int64(15), *storedMetric.Delta)
}

func TestLoad(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	// Prepare initial data
	metrics := []models.Metric{
		{Name: "gauge_metric", Type: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()},
		{Name: "counter_metric", Type: models.CounterType, Delta: func() *int64 { v := int64(10); return &v }()},
	}

	file, err := os.OpenFile(tempFile.Name(), os.O_WRONLY, 0644)
	assert.NoError(t, err)
	defer file.Close()

	// Write each metric as a JSON line in the .txt file
	for _, metric := range metrics {
		data, err := json.Marshal(metric)
		assert.NoError(t, err)
		file.Write(data)
		file.Write([]byte("\n"))
	}

	storage := New(tempFile.Name())
	err = storage.Load(tempFile.Name())
	assert.NoError(t, err)

	loadedMetrics, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, *loadedMetrics, 2)
}

func TestSave(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())

	// Add metrics to the storage
	storage.SetGauge(dto.PostMetricDto{ID: "gauge_metric", MType: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()})
	storage.SetCounter(dto.PostMetricDto{ID: "counter_metric", MType: models.CounterType, Delta: func() *int64 { v := int64(10); return &v }()})

	// Save metrics to file
	err := storage.Save(tempFile.Name())
	assert.NoError(t, err)

	// Verify file content by reloading it
	file, err := os.Open(tempFile.Name())
	assert.NoError(t, err)
	defer file.Close()

	var loadedMetrics []models.Metric
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric models.Metric
		err := json.Unmarshal(scanner.Bytes(), &metric)
		assert.NoError(t, err)
		loadedMetrics = append(loadedMetrics, metric)
	}

	assert.Len(t, loadedMetrics, 2)
	assert.Equal(t, "gauge_metric", loadedMetrics[0].Name)
	assert.Equal(t, "counter_metric", loadedMetrics[1].Name)
}
