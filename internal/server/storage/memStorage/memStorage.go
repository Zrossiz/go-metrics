package memstorage

import "github.com/Zrossiz/go-metrics/internal/server/storage"

type MemStorage struct {
	Metrics []storage.Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: []storage.Metric{},
	}
}

func (m *MemStorage) SetGauge(name string, value float64) bool {
	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			m.Metrics[i].Value = value
			return true
		}
	}

	m.Metrics = append(m.Metrics, storage.Metric{
		Type:  storage.GaugeType,
		Name:  name,
		Value: value,
	})

	return true
}

func (m *MemStorage) SetCounter(name string, value int64) bool {
	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			currentValue, ok := m.Metrics[i].Value.(int64)
			if !ok {
				return false
			}
			m.Metrics[i].Value = currentValue + value
			return true
		}
	}

	m.Metrics = append(m.Metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  name,
		Value: value,
	})

	return true
}

func (m *MemStorage) GetMetric(name string) storage.Metric {
	for i := 0; i < len(m.Metrics); i++ {
		if m.Metrics[i].Name == name {
			return m.Metrics[i]
		}
	}
	return storage.Metric{}
}
