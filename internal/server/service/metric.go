// Package service provides the business logic layer for managing metrics.
package service

import (
	"fmt"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
)

// MetricService handles the business logic for managing metrics.
type MetricService struct {
	storage Storager
}

// Storager defines the interface for storage operations related to metrics.
type Storager interface {
	// SetGauge saves or updates a gauge metric in the storage.
	SetGauge(body dto.PostMetricDto) error
	// SetCounter saves or updates a counter metric in the storage.
	SetCounter(body dto.PostMetricDto) error
	// SetBatch performs a batch update of metrics in the storage.
	SetBatch(body []dto.PostMetricDto) error
	// Get retrieves a single metric by name from the storage.
	Get(name string) (*models.Metric, error)
	// GetAll retrieves all metrics from the storage.
	GetAll() (*[]models.Metric, error)
	// Ping checks the connectivity of the storage.
	Ping() error
}

// New creates a new MetricService with the given storage.
func New(stor Storager) *MetricService {
	return &MetricService{
		storage: stor,
	}
}

// Create processes a single metric DTO and saves it to the storage.
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

// SetBatch processes and saves multiple metrics at once.
func (m *MetricService) SetBatch(body []dto.PostMetricDto) error {
	counterMap := make(map[string]int64)
	newBody := make([]dto.PostMetricDto, 0, len(body))

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

// Get retrieves a single metric by name from the storage.
func (m *MetricService) Get(name string) (*models.Metric, error) {
	metric, err := m.storage.Get(name)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

// GetStringValueMetric retrieves a metric's value as a string by name.
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

// GetAll retrieves all metrics from the storage.
func (m *MetricService) GetAll() (*[]models.Metric, error) {
	metrics, err := m.storage.GetAll()
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

// PingDB checks the connectivity of the storage and returns an error if unavailable.
func (m *MetricService) PingDB() error {
	return m.storage.Ping()
}
