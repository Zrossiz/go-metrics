package filestorage

import (
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"golang.org/x/exp/rand"
)

type FileStorage struct {
	data []models.Metric
	mu   sync.Mutex
}

func New() *FileStorage {
	return &FileStorage{
		data: make([]models.Metric, 0),
	}
}

func (f *FileStorage) SetGauge(metric dto.PostMetricDto) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if metric.Name == f.data[i].Name {
			f.data[i].Value = metric.Value
			return nil
		}
	}

	f.data = append(f.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.Name,
		Type:  models.GaugeType,
		Value: metric.Value,
	})

	return nil
}

func (f *FileStorage) SetCounter(metric dto.PostMetricDto) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if metric.Name == f.data[i].Name {
			f.data[i].Delta += int64(metric.Value)
			return nil
		}
	}

	f.data = append(f.data, models.Metric{
		ID:    uint(rand.Int63()),
		Name:  metric.Name,
		Type:  models.CounterType,
		Delta: int64(metric.Value),
	})

	return nil
}

func (f *FileStorage) Get(body dto.GetMetricDto) (*models.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < len(f.data); i++ {
		if f.data[i].Name == body.Name {
			return &f.data[i], nil
		}
	}

	return nil, nil
}

func (f *FileStorage) GetAll() (*[]models.Metric, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return &f.data, nil
}
