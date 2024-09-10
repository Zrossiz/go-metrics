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

func (d *DBStorage) SetGauge(metric dto.PostMetricDto) error {
	DBMetric := models.Metric{
		Name:  metric.Name,
		Type:  metric.MType,
		Value: metric.Value,
	}

	err := d.db.Create(&DBMetric).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) SetCounter(metric dto.PostMetricDto) error {
	DBMetric := models.Metric{
		Name:  metric.Name,
		Type:  metric.MType,
		Delta: int64(metric.Value),
	}

	err := d.db.Create(&DBMetric).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Get(body dto.GetMetricDto) (*models.Metric, error) {
	var metric models.Metric
	err := d.db.Where("name = ?", body.Name).Last(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (d *DBStorage) GetAll() (*[]models.Metric, error) {
	var metrics []models.Metric

	err := d.db.Order("created_at DESC").Limit(27).Find(&metrics).Error
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}
