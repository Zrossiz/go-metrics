package collector

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
)

func CollectMetrics() []types.Metric {
	var metrics []types.Metric
	seed := time.Now().UnixNano()
	localRand := rand.New(rand.NewSource(seed))
	randomValue := localRand.Float64()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "Alloc",
		Value: float64(memStats.Alloc),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "BuckHashSys",
		Value: float64(memStats.BuckHashSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "Frees",
		Value: float64(memStats.Frees),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "GCCPUFraction",
		Value: float64(memStats.GCCPUFraction),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "GCSys",
		Value: float64(memStats.GCSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapAlloc",
		Value: float64(memStats.HeapAlloc),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapIdle",
		Value: float64(memStats.HeapIdle),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapInuse",
		Value: float64(memStats.HeapInuse),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapObjects",
		Value: float64(memStats.HeapObjects),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapReleased",
		Value: float64(memStats.HeapReleased),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "HeapSys",
		Value: float64(memStats.HeapSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "LastGC",
		Value: float64(memStats.LastGC),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "Lookups",
		Value: float64(memStats.Lookups),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "MCacheInuse",
		Value: float64(memStats.MCacheInuse),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "MCacheSys",
		Value: float64(memStats.MCacheSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "Mallocs",
		Value: float64(memStats.Mallocs),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "NextGC",
		Value: float64(memStats.NextGC),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "NumForcedGC",
		Value: float64(memStats.NumForcedGC),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "NumGC",
		Value: float64(memStats.NumGC),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "OtherSys",
		Value: float64(memStats.OtherSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "PauseTotalNs",
		Value: float64(memStats.PauseTotalNs),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "StackInuse",
		Value: float64(memStats.StackInuse),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "StackSys",
		Value: float64(memStats.StackSys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "Sys",
		Value: float64(memStats.Sys),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "TotalAlloc",
		Value: float64(memStats.TotalAlloc),
	})
	metrics = append(metrics, types.Metric{
		Type:  constants.Gauge,
		Name:  "RandomValue",
		Value: randomValue,
	})

	return metrics
}
