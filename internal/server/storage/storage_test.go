package storage_test

import (
	"context"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/filestorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/memstorage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock for database connection pool.
type MockDBPool struct {
	mock.Mock
}

func TestNew_WithDBStorage(t *testing.T) {
	db, _ := pgxpool.Connect(context.Background(), "postgresql://postgres:root@localhost/metrics")
	logger := zap.NewNop()
	cfg := &config.Config{DBDSN: "mock_dsn"}

	st := storage.New(db, cfg, logger)

	assert.IsType(t, &dbstorage.DBStorage{}, st, "Expected DBStorage when DSN is provided")
}

// func setupDB(t *testing.T) (*pgxpool.Pool, *zap.Logger) {
// 	db, err := pgxpool.Connect(context.Background(), "postgresql://postgres:root@localhost/metrics")
// 	require.NoError(t, err, "Database connection failed")
// 	logger := zap.NewNop()
// 	return db, logger
// }

func TestNew_WithFileStorage(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		FileStoragePath: "/tmp/metrics.json",
		StoreInterval:   "10s",
		Restore:         true,
	}

	st := storage.New(nil, cfg, logger)

	assert.IsType(t, &filestorage.FileStorage{}, st, "Expected FileStorage when FileStoragePath is provided")
}

func TestNew_WithMemStorage(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{}

	st := storage.New(nil, cfg, logger)

	assert.IsType(t, &memstorage.MemStorage{}, st, "Expected MemStorage when no other storage is configured")
}

func TestStorage_SetBatch(t *testing.T) {
	st := memstorage.New()

	metrics := []dto.PostMetricDto{
		{
			ID:    "test_gauge1",
			MType: "gauge",
			Value: func() *float64 { v := 1.23; return &v }(),
		},
		{
			ID:    "test_counter",
			MType: "counter",
			Delta: func() *int64 { d := int64(10); return &d }(),
		},
	}

	err := st.SetBatch(metrics)
	assert.NoError(t, err, "Setting a batch of metrics should not return an error")

	gauge, err := st.Get("test_gauge1")
	assert.NoError(t, err, "Getting a gauge metric should not return an error")
	assert.Equal(t, *metrics[0].Value, *gauge.Value, "Gauge value should match")

	counter, err := st.Get("test_counter")
	assert.NoError(t, err, "Getting a counter metric should not return an error")
	assert.Equal(t, *metrics[1].Delta, *counter.Delta, "Counter delta should match")
}

func TestStorage_SaveAndLoad(t *testing.T) {
	st := filestorage.New("/tmp/metrics_test.json")

	metrics := []dto.PostMetricDto{
		{
			ID:    "test_gauge",
			MType: "gauge",
			Value: func() *float64 { v := 1.23; return &v }(),
		},
	}
	err := st.SetBatch(metrics)
	assert.NoError(t, err, "Setting metrics should not return an error")

	err = st.Save("/tmp/metrics_test.json")
	assert.NoError(t, err, "Saving metrics should not return an error")

	// Create a new storage instance to simulate loading
	newStorage := filestorage.New("/tmp/metrics_test.json")
	err = newStorage.Load("/tmp/metrics_test.json")
	assert.NoError(t, err, "Loading metrics should not return an error")

	retrieved, err := newStorage.Get("test_gauge")
	assert.NoError(t, err, "Getting a loaded metric should not return an error")
	assert.Equal(t, *metrics[0].Value, *retrieved.Value, "Loaded metric value should match")
}
