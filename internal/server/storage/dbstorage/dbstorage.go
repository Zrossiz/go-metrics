package dbstorage

import (
	"fmt"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBStorage struct {
	db     *gorm.DB
	logger *zap.Logger
}

const maxRetries = 3
const retryDelay = 1 * time.Second

func New(dbConn *gorm.DB, log *zap.Logger) *DBStorage {
	return &DBStorage{db: dbConn, logger: log}
}

func GetConnect(connStr string, logger *zap.Logger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	delay := retryDelay

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err == nil {
			break
		}

		logger.Error("error connect to db. retry...")
		delay += 2 * time.Second
		time.Sleep(delay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
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
		Type:  metric.Type,
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
		Type:  metric.Type,
		Delta: int64(metric.Value),
	}

	err := d.db.Create(&DBMetric).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Get(name string) (*models.Metric, error) {
	var metric models.Metric
	err := d.db.Where("name = ?", name).Last(&metric).Error
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

	pollCount, _ := d.Get("PollCount")

	metrics = append(metrics, *pollCount)

	return &metrics, nil
}

func (d *DBStorage) SetBatch(body []dto.PostMetricDto) error {
	var metrics []models.Metric

	for i := 0; i < len(body); i++ {
		if body[i].Type == models.CounterType {
			pollCount, _ := d.Get(body[i].Name)

			if pollCount != nil {
				d.db.Model(&pollCount).Update("delta", pollCount.Delta+int64(body[i].Value))
				continue
			}

			metrics = append(metrics, models.Metric{
				Name:  body[i].Name,
				Type:  body[i].Type,
				Delta: int64(body[i].Value),
			})
			continue
		}

		metrics = append(metrics, models.Metric{
			Name:  body[i].Name,
			Type:  body[i].Type,
			Value: body[i].Value,
		})
	}

	err := d.db.Create(&metrics).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Load(filePath string) error {
	return nil
}

func (d *DBStorage) Save(filePath string) error {
	return nil
}

func (d *DBStorage) Close(string) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func MigrateSQL(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Metric{})
	if err != nil {
		return err
	}

	return nil
}
