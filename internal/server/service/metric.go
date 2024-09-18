package service

import (
	"fmt"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type MetricService struct {
	storage Storager
	logger  *zap.Logger
	dbConn  *pgxpool.Pool
}

type Storager interface {
	SetGauge(body dto.PostMetricDto) error
	SetCounter(body dto.PostMetricDto) error
	SetBatch(body []dto.PostMetricDto) error
	Get(name string) (*models.Metric, error)
	GetAll() (*[]models.Metric, error)
	Ping() error
}

func New(stor Storager, logger *zap.Logger, dbConn *pgxpool.Pool) *MetricService {
	return &MetricService{
		storage: stor,
		logger:  logger,
	}
}

func (m *MetricService) Create(body dto.PostMetricDto) error {
	if body.MType == models.CounterType {
		if body.Delta == nil {
			return fmt.Errorf("delta not found")
		}

		err := m.storage.SetCounter(body)
		if err != nil {
			return err
		}

		return nil
	}

	if body.Value == nil {
		return fmt.Errorf("value not found")
	}

	err := m.storage.SetGauge(body)
	if err != nil {
		return err
	}

	return nil
}

func (m *MetricService) SetBatch(body []dto.PostMetricDto) error {
	counterMap := make(map[string]int64)
	var newBody []dto.PostMetricDto

	for _, metric := range body {
		if metric.MType == models.CounterType {
			counterMap[metric.ID] += *metric.Delta
		} else {
			newBody = append(newBody, metric)
		}
	}

	for id, value := range counterMap {
		newBody = append(newBody, dto.PostMetricDto{
			ID:    id,
			MType: models.CounterType,
			Delta: &value,
		})
	}

	err := m.storage.SetBatch(newBody)
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
			value = fmt.Sprintf("%v", *metric.Delta)
		} else {
			value = fmt.Sprintf("%v", *metric.Value)
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
	return m.storage.Ping()
}
