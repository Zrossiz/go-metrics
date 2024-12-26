package dbstorage

import (
	"context"

	"github.com/Zrossiz/go-metrics/internal/server/models"
)

type DataBase interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, name string) (*models.Metric, error)
	GetAll(ctx context.Context) ([]models.Metric, error)
	SetGauge(ctx context.Context, metric models.Metric) error
	SetCounter(ctx context.Context, metric models.Metric) error
	SetBatch(ctx context.Context, metrics []models.Metric) error
	Close() error
}
