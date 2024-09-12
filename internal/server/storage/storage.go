package storage

import (
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
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
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
	Load(filePath string) error
	Save(filePath string) error
}

func New(dbConn *gorm.DB, cfg *config.Config, log *zap.Logger) Storage {
	if dbConn != nil {
		return dbstorage.New(dbConn)
	}

	if len(cfg.FileStoragePath) > 0 {
		store := filestorage.New(cfg.FileStoragePath)

		if cfg.Restore {
			log.Info("start collect metrics from storage...")
			err := store.Load(cfg.FileStoragePath)
			if err != nil {
				log.Fatal("error collect metric", zap.Error(err))
			}
			log.Info("metric collected")
		}

		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		defer func() {
			log.Info("Stopping ticker")
			ticker.Stop()
		}()
		stop := make(chan bool)

		go func() {
			for {
				select {
				case <-ticker.C:
					log.Info("Saving metrics to file", zap.String("file", cfg.FileStoragePath))
					if err := store.Save(cfg.FileStoragePath); err != nil {
						log.Error("Failed to save metrics to file", zap.Error(err))
					}
					log.Info("Successful save")
				case <-stop:
					log.Info("Stopping task execution")
				}
			}
		}()
	}

	return memstorage.New()
}
