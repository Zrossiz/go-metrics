package storage

import (
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/filestorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/memstorage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	Get(body dto.GetMetricDto) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
}

func New(dbConn *gorm.DB, filePath string, storeInterval time.Duration, log *zap.Logger) Storage {
	if dbConn != nil {
		return dbstorage.New(dbConn, log)
	}

	if storeInterval > 0 {
		return filestorage.New()
	}

	return memstorage.New()
}
