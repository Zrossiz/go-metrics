package memstorage

import "github.com/Zrossiz/go-metrics/internal/server/storage"

type MemStorage struct {
	metrics []storage.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: []storage.Metric{},
	}
}

func (m *MemStorage) SetGauge(name string, value float64) bool {
	for i := 0; i < len(m.metrics); i++ {
		if m.metrics[i].Name == name {
			m.metrics[i].Value = value
			return true
		}
	}

	m.metrics = append(m.metrics, storage.Metric{
		Type:  storage.GaugeType,
		Name:  name,
		Value: value,
	})

	return true
}

func (m *MemStorage) SetCounter(name string, value int64) bool {
	for i := 0; i < len(m.metrics); i++ {
		if m.metrics[i].Name == name {
			currentValue, ok := m.metrics[i].Value.(int64)
			if !ok {
				return false
			}
			m.metrics[i].Value = currentValue + value
			return true
		}
	}

	m.metrics = append(m.metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  name,
		Value: value,
	})

	return true
}

func (m *MemStorage) GetMetric(name string) storage.Metric {
	for i := 0; i < len(m.metrics); i++ {
		if m.metrics[i].Name == name {
			return m.metrics[i]
		}
	}
	return storage.Metric{}
}
