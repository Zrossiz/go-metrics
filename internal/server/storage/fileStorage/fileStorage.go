package filestorage

import (
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
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

func (f *FileStorage) CreateGauge(metric dto.PostMetricDto) error {
	return nil
}

func (f *FileStorage) CreateCounter(metric dto.PostMetricDto) error {
	return nil
}

func (f *FileStorage) Get(body dto.GetMetricDto) (models.Metric, error) {
	return models.Metric{}, nil
}
