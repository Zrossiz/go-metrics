package collector

import (
	"math/rand"
	"reflect"
	"runtime"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func CollectMetrics() []types.Metric {
	var metrics []types.Metric
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	v := reflect.ValueOf(memStats)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i)

		var value float64

		switch fieldValue.Kind() {
		case reflect.Float32, reflect.Float64:
			value = fieldValue.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = float64(fieldValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = float64(fieldValue.Uint())
		default:
			continue
		}

		metrics = append(metrics, types.Metric{
			Type:  constants.Gauge,
			Name:  fieldName,
			Value: value,
		})
	}

	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "RandomValue",
		Value: rand.Float64(),
	})

	return metrics
}
