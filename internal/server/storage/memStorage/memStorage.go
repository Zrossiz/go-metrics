package memstorage

import (
	"sync"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
)

type MemStorage struct {
	Metrics []storage.Metric
}

var MemStore *MemStorage

func NewMemStorage() {
	MemStore = &MemStorage{
		Metrics: []storage.Metric{},
	}
}

func (m *MemStorage) SetGauge(name string, value float64) *storage.Metric {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			m.Metrics[i].Value = value
			return &m.Metrics[i]
		}
	}

	newMetric := storage.Metric{
		Type:  storage.GaugeType,
		Name:  name,
		Value: value,
	}

	m.Metrics = append(m.Metrics, newMetric)

	return &newMetric
}

func (m *MemStorage) SetCounter(name string, value int64) *storage.Metric {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			currentValue, ok := m.Metrics[i].Value.(int64)
			if !ok {
				return nil
			}
			m.Metrics[i].Value = currentValue + value
			return &m.Metrics[i]
		}
	}

	newMetric := storage.Metric{
		Type:  storage.CounterType,
		Name:  name,
		Value: value,
	}

	m.Metrics = append(m.Metrics, newMetric)

	return &newMetric
}

func (m *MemStorage) GetMetric(name string) *storage.Metric {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			return &m.Metrics[i]
		}
	}
	return nil
}
