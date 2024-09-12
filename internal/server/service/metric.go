package service

import (
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"go.uber.org/zap"
)

type MetricService struct {
	storage Storager
	logger  *zap.Logger
}

type Storager interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
}

func New(stor Storager, logger *zap.Logger) *MetricService {
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

func (m *MetricService) Get(name string) (*models.Metric, error) {
	metric, err := m.storage.Get(name)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

func (m *MetricService) GetStringValueMetric(name string) (string, error) {
	return "", nil
}

func (m *MetricService) GetAll() (*[]models.Metric, error) {
	return &[]models.Metric{}, nil
}
