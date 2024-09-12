package storage

import (
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/filestorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/memstorage"
	"gorm.io/gorm"
)

type Storage interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
	Load() error
	Save() error
}

func New(dbConn *gorm.DB, filePath string) Storage {
	if dbConn != nil {
		return dbstorage.New(dbConn)
	}

	if len(filePath) > 0 {
		return filestorage.New(filePath)
	}

	return memstorage.New()
}
