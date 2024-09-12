package memstorage

import (
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"golang.org/x/exp/rand"
)

type MemStorage struct {
	data []models.Metric
	mu   sync.Mutex
}

func New() *MemStorage {
	return &MemStorage{
		data: make([]models.Metric, 0),
	}
}

func (m *MemStorage) SetGauge(metric dto.PostMetricDto) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.data); i++ {
		if metric.Name == m.data[i].Name {
			m.data[i].Value = metric.Value
			return nil
		}
	}

	m.data = append(m.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.Name,
		Type:  models.GaugeType,
		Value: metric.Value,
	})

	return nil
}

func (m *MemStorage) SetCounter(metric dto.PostMetricDto) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.data); i++ {
		if metric.Name == m.data[i].Name {
			m.data[i].Delta += int64(metric.Value)
			return nil
		}
	}

	m.data = append(m.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.Name,
		Type:  models.CounterType,
		Delta: int64(metric.Value),
	})

	return nil
}

func (m *MemStorage) Get(name string) (*models.Metric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.data); i++ {
		if m.data[i].Name == name {
			return &m.data[i], nil
		}
	}

	return nil, nil
}

func (m *MemStorage) GetAll() (*[]models.Metric, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return &m.data, nil
}

func (m *MemStorage) Load() error {
	return nil
}

func (m *MemStorage) Save() error {
	return nil
}
