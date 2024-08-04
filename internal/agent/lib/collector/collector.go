package collector

import (
	"math/rand"
	"runtime"
	"time"
)

type MetricType string

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

func CollectMetrics() []Metric {
	var metrics []Metric
	seed := time.Now().UnixNano()
    localRand := rand.New(rand.NewSource(seed))
    randomValue := localRand.Float64()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "Alloc",
		Value: float64(memStats.Alloc),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "BuckHashSys",
		Value: float64(memStats.BuckHashSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "Frees",
		Value: float64(memStats.Frees),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "GCCPUFraction",
		Value: float64(memStats.GCCPUFraction),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "GCSys",
		Value: float64(memStats.GCSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapAlloc",
		Value: float64(memStats.HeapAlloc),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapIdle",
		Value: float64(memStats.HeapIdle),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapInuse",
		Value: float64(memStats.HeapInuse),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapObjects",
		Value: float64(memStats.HeapObjects),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapReleased",
		Value: float64(memStats.HeapReleased),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "HeapSys",
		Value: float64(memStats.HeapSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "LastGC",
		Value: float64(memStats.LastGC),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "Lookups",
		Value: float64(memStats.Lookups),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "MCacheInuse",
		Value: float64(memStats.MCacheInuse),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "MCacheSys",
		Value: float64(memStats.MCacheSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "Mallocs",
		Value: float64(memStats.Mallocs),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "NextGC",
		Value: float64(memStats.NextGC),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "NumForcedGC",
		Value: float64(memStats.NumForcedGC),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "NumGC",
		Value: float64(memStats.NumGC),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "OtherSys",
		Value: float64(memStats.OtherSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "PauseTotalNs",
		Value: float64(memStats.PauseTotalNs),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "StackInuse",
		Value: float64(memStats.StackInuse),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "StackSys",
		Value: float64(memStats.StackSys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "Sys",
		Value: float64(memStats.Sys),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "TotalAlloc",
		Value: float64(memStats.TotalAlloc),
	})
	metrics = append(metrics, Metric{
		Type:  Gauge,
		Name:  "RandomValue",
		Value: randomValue,
	})

	return metrics
}
