package memstorage

import (
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
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

func (m *MemStorage) CreateGauge(metric dto.PostMetricDto) error {
	return nil
}

func (m *MemStorage) CreateCounter(metric dto.PostMetricDto) error {
	return nil
}

func (m *MemStorage) Get(body dto.GetMetricDto) (models.Metric, error) {
	return models.Metric{}, nil
}
