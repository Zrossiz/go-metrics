package memstorage

import (
	storageServer "github.com/Zrossiz/go-metrics/internal/storage/server"
)

var Metrics []storageServer.Metric

func SetGauge(name string, value float64) bool {
	for i := 0; i < len(Metrics); i++ {
		if Metrics[i].Name == name {
			Metrics[i].Value = value
			return true
		}
	}

	Metrics = append(Metrics, storageServer.Metric{
		Type:  storageServer.GaugeType,
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

	Metrics = append(Metrics, storageServer.Metric{
		Type:  storageServer.CounterType,
		Name:  name,
		Value: value,
	})

	return true
}

func GetMetric(name string) storageServer.Metric {
	for i := 0; i < len(Metrics); i++ {
		if Metrics[i].Name == name {
			return Metrics[i]
		}
	}
	return storageServer.Metric{}
}
