package dbstorage

import (
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBStorage struct {
	db  *gorm.DB
	log *zap.Logger
}

func New(dbConn *gorm.DB, log *zap.Logger) *DBStorage {
	return &DBStorage{db: dbConn, log: log}
}

func GetConnect(connStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = Ping(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) CreateGauge(metric dto.PostMetricDto) error {
	return nil
}

func (d *DBStorage) CreateCounter(metric dto.PostMetricDto) error {
	return nil
}

func (d *DBStorage) Get(body dto.GetMetricDto) (models.Metric, error) {
	return models.Metric{}, nil
}
