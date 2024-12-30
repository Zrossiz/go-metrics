// Package storage provides a unified interface for managing metrics in various storage backends,
// including memory, file, and database-based storages.
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

// Storage defines the interface for metric storage operations.
// It supports multiple implementations for handling metrics: in memory, files, or databases.
type Storage interface {
	// SetGauge saves or updates a gauge metric in the storage.
	SetGauge(body dto.PostMetricDto) error
	// SetCounter saves or updates a counter metric in the storage.
	SetCounter(body dto.PostMetricDto) error
	// SetBatch performs a batch update of metrics in the storage.
	SetBatch(body []dto.PostMetricDto) error
	// Get retrieves a single metric by name from the storage.
	Get(name string) (*models.Metric, error)
	// GetAll retrieves all metrics from the storage.
	GetAll() ([]models.Metric, error)
	// Load loads metrics from a file into the storage.
	Load(filePath string) error
	// Save persists the current metrics from the storage to a file.
	Save(filePath string) error
	// Close performs any cleanup or connection closure for the storage.
	Close() error
	// Ping checks the connectivity of the storage, if applicable.
	Ping() error
}

// New initializes and returns an appropriate storage implementation based on the provided configuration.
// Priority of storage selection is as follows:
// 1. Database storage if a database connection string is provided.
// 2. File storage if a file path is configured.
// 3. In-memory storage if neither database nor file storage is configured.
//
// Parameters:
// - dbConn: A connection pool for interacting with the database.
// - cfg: Configuration object containing storage-related settings.
// - log: Logger instance for logging storage operations.
func New(dbConn *pgxpool.Pool, cfg *config.Config, log *zap.Logger) Storage {
	// Use database storage if a DSN is provided.
	if cfg.DBDSN != "" {
		fmt.Println("selected: db storage")
		return dbstorage.New(dbConn, log)
	}

	// Use file storage if a file storage path is provided.
	if len(cfg.FileStoragePath) > 0 {
		fmt.Println("selected: file storage")
		store := filestorage.New(cfg.FileStoragePath)

		// Restore metrics from the file if configured to do so.
		if cfg.Restore {
			log.Info("start collect metrics from storage...")
			err := store.Load(cfg.FileStoragePath)
			if err != nil {
				log.Fatal("error collect metric", zap.Error(err))
			}
			log.Info("metric collected")
		}

		// Schedule periodic saving of metrics to the file.
		storeIntervalDuration, err := time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			fmt.Println("invalid duration")
		}

		ticker := time.NewTicker(storeIntervalDuration)
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

	// Default to in-memory storage.
	fmt.Println("selected: mem storage")
	return memstorage.New()
}
