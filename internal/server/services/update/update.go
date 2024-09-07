package update

import (
	"strconv"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
)

func JSONMetric(body dto.MetricDTO) *dto.MetricDTO {
	var updatedMetric *storage.Metric

	switch body.MType {
	case storage.GaugeType:
		if body.Value == nil {
			return nil
		}
		updatedMetric = memstorage.MemStore.SetGauge(body.ID, *body.Value)
	case storage.CounterType:
		if body.Delta == nil {
			return nil
		}
		updatedMetric = memstorage.MemStore.SetCounter(body.ID, *body.Delta)
	default:
		return nil
	}

	if updatedMetric == nil {
		return nil
	}

	responseMetric := dto.MetricDTO{
		ID:    updatedMetric.Name,
		MType: updatedMetric.Type,
	}

	if v, ok := updatedMetric.Value.(float64); ok {
		responseMetric.Value = &v
	}

	if d, ok := updatedMetric.Value.(int64); ok {
		responseMetric.Delta = &d
	}

	return &responseMetric
}

func ParamMetric(typeMetric string, nameMetric string, valueMetric string) (*dto.MetricDTO, error) {

	responseMetric := dto.MetricDTO{
		ID:    nameMetric,
		MType: typeMetric,
	}

	switch typeMetric {
	case storage.GaugeType:
		float64MetricValue, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			return nil, err
		}

		memstorage.MemStore.SetGauge(nameMetric, float64MetricValue)
		responseMetric.Value = &float64MetricValue
	case storage.CounterType:
		int64MetricValue, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			return nil, err
		}

		memstorage.MemStore.SetCounter(nameMetric, int64MetricValue)
		responseMetric.Delta = &int64MetricValue
	default:
	}

	return &responseMetric, nil
}
