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

func TestSetGauge_ExistingMetric(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())

	metric := dto.PostMetricDto{ID: "gauge_metric", MType: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()}
	err := storage.SetGauge(metric)
	assert.NoError(t, err)

	updatedValue := 678.90
	metric.Value = &updatedValue
	err = storage.SetGauge(metric)
	assert.NoError(t, err)

	storedMetric, err := storage.Get("gauge_metric")
	assert.NoError(t, err)
	assert.Equal(t, &updatedValue, storedMetric.Value)
}

func TestSetCounter_ExistingMetric(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	initialDelta := int64(10)
	metric := dto.PostMetricDto{ID: "counter_metric", MType: models.CounterType, Delta: &initialDelta}

	err := storage.SetCounter(metric)
	assert.NoError(t, err)

	additionalDelta := int64(5)
	metric.Delta = &additionalDelta
	err = storage.SetCounter(metric)
	assert.NoError(t, err)

	storedMetric, err := storage.Get("counter_metric")
	assert.NoError(t, err)
	assert.Equal(t, int64(15), *storedMetric.Delta)
}

func TestSetCounter_ZeroDelta(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	metric := dto.PostMetricDto{ID: "counter_metric_zero", MType: models.CounterType, Delta: func() *int64 { v := int64(0); return &v }()}

	err := storage.SetCounter(metric)
	assert.NoError(t, err)

	storedMetric, err := storage.Get("counter_metric_zero")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), *storedMetric.Delta)
}

func TestLoad_EmptyFile(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	err := storage.Load(tempFile.Name())
	assert.NoError(t, err)

	// Проверяем, что хранилище пустое
	loadedMetrics, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, loadedMetrics, 0)
}

func TestLoad_InvalidData(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	file, err := os.OpenFile(tempFile.Name(), os.O_WRONLY, 0644)
	assert.NoError(t, err)
	defer file.Close()

	_, err = file.WriteString("invalid data\n")
	assert.NoError(t, err)

	storage := New(tempFile.Name())
	err = storage.Load(tempFile.Name())
	assert.NoError(t, err)

	loadedMetrics, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, loadedMetrics, 0)
}

func TestSave_EmptyStorage(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	err := storage.Save(tempFile.Name())
	assert.NoError(t, err)

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

	assert.Len(t, loadedMetrics, 0)
}

func TestSave_CorrectData(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	storage.SetGauge(dto.PostMetricDto{ID: "gauge_metric", MType: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()})
	storage.SetCounter(dto.PostMetricDto{ID: "counter_metric", MType: models.CounterType, Delta: func() *int64 { v := int64(10); return &v }()})

	err := storage.Save(tempFile.Name())
	assert.NoError(t, err)

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
	assert.Len(t, loadedMetrics, 2)
}

func TestSetBatch(t *testing.T) {
	tempFile := createTempFile(t)
	defer os.Remove(tempFile.Name())

	storage := New(tempFile.Name())
	batch := []dto.PostMetricDto{
		{ID: "gauge_metric_batch", MType: models.GaugeType, Value: func() *float64 { v := 123.45; return &v }()},
		{ID: "counter_metric_batch", MType: models.CounterType, Delta: func() *int64 { v := int64(10); return &v }()},
	}

	err := storage.SetBatch(batch)
	assert.NoError(t, err)

	loadedMetrics, err := storage.GetAll()
	assert.NoError(t, err)
	assert.Len(t, loadedMetrics, 2)
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
