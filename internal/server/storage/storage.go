package storage

import (
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/filestorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/memstorage"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	SetBatch(body []dto.PostMetricDto) error
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
	Load(filePath string) error
	Save(filePath string) error
	Close(filePath string) error
}

func New(dbConn *pgxpool.Pool, cfg *config.Config, log *zap.Logger) Storage {
	if cfg.DBDSN != "" {
		fmt.Println("selected: db storage")
		return dbstorage.New(dbConn, log)
	}

	if len(cfg.FileStoragePath) > 0 {
		fmt.Println("selected: file storage")
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
					return
				}
			}
		}()

		return store
	}

	fmt.Println("selected: mem storage")
	return memstorage.New()

}
