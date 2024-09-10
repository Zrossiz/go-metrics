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
	CreateGauge(body dto.PostMetricDto) error
	CreateCounter(body dto.PostMetricDto) error
	Get(body dto.GetMetricDto) (models.Metric, error)
	GetHTML()
}

func New(stor Storager, logger *zap.Logger) *MetricService {
	return &MetricService{
		storage: stor,
		logger:  logger,
	}
}

func (m *MetricService) GreateGauge(body dto.PostMetricDto) error {
	return nil
}

func (m *MetricService) CreateCounter(body dto.PostMetricDto) error {
	return nil
}

func (m *MetricService) Get() (models.Metric, error) {
	return models.Metric{}, nil
}
