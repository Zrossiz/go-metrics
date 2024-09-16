package service

import (
	"fmt"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MetricService struct {
	storage Storager
	logger  *zap.Logger
	dbConn  *gorm.DB
}

type Storager interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	SetBatch(body []dto.PostMetricDto) error
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
}

func New(stor Storager, logger *zap.Logger, dbConn *gorm.DB) *MetricService {
	return &MetricService{
		storage: stor,
		logger:  logger,
	}
}

func (m *MetricService) Create(body dto.PostMetricDto) error {
	if body.MType == models.CounterType {
		err := m.storage.SetCounter(body)
		if err != nil {
			return err
		}

		return nil
	}

	err := m.storage.SetGauge(body)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricService) SetBatch(body []dto.PostMetricDto) error {
	err := m.storage.SetBatch(body)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricService) Get(name string) (*models.Metric, error) {
	metric, err := m.storage.Get(name)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

func (m *MetricService) GetStringValueMetric(name string) (string, error) {
	metric, err := m.storage.Get(name)
	if err != nil {
		return "", err
	}

	if metric != nil {
		var value string
		if metric.Type == models.CounterType {
			value = fmt.Sprintf("%v", metric.Delta)
		} else {
			value = fmt.Sprintf("%v", metric.Value)
		}

		return value, nil
	}

	return "", nil
}

func (m *MetricService) GetAll() (*[]models.Metric, error) {
	metrics, err := m.storage.GetAll()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (m *MetricService) PingDB() error {
	return dbstorage.Ping(m.dbConn)
}
