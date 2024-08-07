package memstorage

import "github.com/Zrossiz/go-metrics/server/internal/storage"

var Metrics []storage.Metric

func SetGauge(name string, value float64) bool {
	for i := 0; i < len(Metrics); i++ {
		if Metrics[i].Name == name {
			Metrics[i].Value = value
			return true
		}
	}

	Metrics = append(Metrics, storage.Metric{
		Type:  storage.GaugeType,
		Name:  name,
		Value: value,
	})

	return true
}

func SetCounter(name string, value int64) bool {
	for i := 0; i < len(Metrics); i++ {
		if Metrics[i].Name == name {
			currentValue, ok := Metrics[i].Value.(int64)
			if !ok {
				return false
			}
			Metrics[i].Value = currentValue + value
			return true
		}
	}

	Metrics = append(Metrics, storage.Metric{
		Type:  storage.CounterType,
		Name:  name,
		Value: value,
	})

	return true
}

func GetMetric(name string) storage.Metric {
	for i := 0; i < len(Metrics); i++ {
		if Metrics[i].Name == name {
			return Metrics[i]
		}
	}
	return storage.Metric{}
}
