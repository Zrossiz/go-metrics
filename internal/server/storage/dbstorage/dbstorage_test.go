package dbstorage

import (
	"context"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

// Mock setup for pgxpool and logger for testing purposes
func setupDB(t *testing.T) (*pgxpool.Pool, *zap.Logger) {
	db, err := pgxpool.Connect(context.Background(), "postgresql://postgres:root@localhost/metrics")
	require.NoError(t, err, "Database connection failed")
	logger := zap.NewNop()
	return db, logger
}

func TestDBStorage_Ping(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	err := storage.Ping()
	assert.NoError(t, err, "Ping should succeed when DB is connected")
}

func TestDBStorage_SetGauge(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	metric := dto.PostMetricDto{
		ID:    "gauge_metric",
		MType: "gauge",
		Value: float64Ptr(42.5),
	}
	err := storage.SetGauge(metric)
	assert.NoError(t, err, "SetGauge should insert gauge metric without error")
}

func TestDBStorage_SetCounter(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	metric := dto.PostMetricDto{
		ID:    "counter_metric",
		MType: "counter",
		Delta: int64Ptr(10),
	}
	err := storage.SetCounter(metric)
	assert.NoError(t, err, "SetCounter should insert or update counter metric without error")
}

func TestDBStorage_Get(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	metric := dto.PostMetricDto{
		ID:    "test_metric",
		MType: "gauge",
		Value: float64Ptr(100.0),
	}
	_ = storage.SetGauge(metric)

	result, err := storage.Get("test_metric")
	assert.NoError(t, err, "Get should retrieve the inserted metric without error")
	assert.NotNil(t, result, "Retrieved metric should not be nil")
	assert.Equal(t, "test_metric", result.Name, "Metric name should match")

	if assert.NotNil(t, result.Value, "Metric value should not be nil") {
		assert.Equal(t, float64(100.0), *result.Value, "Metric value should match")
	}
}

func TestDBStorage_GetAll(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	// Insert multiple test metrics
	metrics := []dto.PostMetricDto{
		{ID: "metric_1", MType: "gauge", Value: float64Ptr(50.0)},
		{ID: "metric_2", MType: "counter", Delta: int64Ptr(5)},
	}
	for _, metric := range metrics {
		_ = storage.SetGauge(metric)
	}

	allMetrics, err := storage.GetAll()
	assert.NoError(t, err, "GetAll should retrieve all inserted metrics without error")
	assert.GreaterOrEqual(t, len(*allMetrics), len(metrics), "Should retrieve at least the metrics we inserted")
}

func TestDBStorage_SetBatch(t *testing.T) {
	db, logger := setupDB(t)
	defer db.Close()
	storage := New(db, logger)

	batchMetrics := []dto.PostMetricDto{
		{ID: "batch_metric_1", MType: "gauge", Value: float64Ptr(150.0)},
		{ID: "batch_metric_2", MType: "counter", Delta: int64Ptr(10)},
	}

	err := storage.SetBatch(batchMetrics)
	assert.NoError(t, err, "SetBatch should insert metrics in batch without error")
}

func TestDBStorage_Close(t *testing.T) {
	db, logger := setupDB(t)
	storage := New(db, logger)

	err := storage.Close()
	assert.NoError(t, err, "Close should close the DB connection without error")

	err = storage.Ping()
	assert.Error(t, err, "Ping should fail after DB is closed")
}

func TestDBStorage_GetConnect_RetryOnError(t *testing.T) {
	logger := zap.NewNop()
	db, err := GetConnect("invalid_conn_str", logger)
	assert.Error(t, err, "GetConnect should fail with invalid connection string")
	assert.Nil(t, db, "Returned db should be nil on failure")
}
