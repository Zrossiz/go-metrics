package collector

import (
	"math/rand"
	"runtime"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetMetrics(counter *int64) []types.Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	*counter += 1

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		vmStat = &mem.VirtualMemoryStat{}
	}

	cpuPercentages, err := cpu.Percent(0, true)
	if err != nil {
		cpuPercentages = []float64{0}
	}

	metrics := []types.Metric{
		{Type: constants.Counter, Name: "PollCount", Value: *counter},
		{Name: "Alloc", Type: "gauge", Value: float64(m.Alloc)},
		{Name: "BuckHashSys", Type: "gauge", Value: float64(m.BuckHashSys)},
		{Name: "Frees", Type: "gauge", Value: float64(m.Frees)},
		{Name: "GCCPUFraction", Type: "gauge", Value: m.GCCPUFraction},
		{Name: "GCSys", Type: "gauge", Value: float64(m.GCSys)},
		{Name: "HeapAlloc", Type: "gauge", Value: float64(m.HeapAlloc)},
		{Name: "HeapIdle", Type: "gauge", Value: float64(m.HeapIdle)},
		{Name: "HeapInuse", Type: "gauge", Value: float64(m.HeapInuse)},
		{Name: "HeapObjects", Type: "gauge", Value: float64(m.HeapObjects)},
		{Name: "HeapReleased", Type: "gauge", Value: float64(m.HeapReleased)},
		{Name: "HeapSys", Type: "gauge", Value: float64(m.HeapSys)},
		{Name: "LastGC", Type: "gauge", Value: float64(m.LastGC)},
		{Name: "Lookups", Type: "gauge", Value: float64(m.Lookups)},
		{Name: "MCacheInuse", Type: "gauge", Value: float64(m.MCacheInuse)},
		{Name: "MCacheSys", Type: "gauge", Value: float64(m.MCacheSys)},
		{Name: "MSpanInuse", Type: "gauge", Value: float64(m.MSpanInuse)},
		{Name: "MSpanSys", Type: "gauge", Value: float64(m.MSpanSys)},
		{Name: "Mallocs", Type: "gauge", Value: float64(m.Mallocs)},
		{Name: "NextGC", Type: "gauge", Value: float64(m.NextGC)},
		{Name: "NumForcedGC", Type: "gauge", Value: float64(m.NumForcedGC)},
		{Name: "NumGC", Type: "gauge", Value: float64(m.NumGC)},
		{Name: "OtherSys", Type: "gauge", Value: float64(m.OtherSys)},
		{Name: "PauseTotalNs", Type: "gauge", Value: float64(m.PauseTotalNs)},
		{Name: "StackInuse", Type: "gauge", Value: float64(m.StackInuse)},
		{Name: "StackSys", Type: "gauge", Value: float64(m.StackSys)},
		{Name: "Sys", Type: "gauge", Value: float64(m.Sys)},
		{Name: "TotalAlloc", Type: "gauge", Value: float64(m.TotalAlloc)},
		{Name: "RandomValue", Type: "gauge", Value: rand.Float64()},
		{Name: "TotalMemory", Type: "gauge", Value: float64(vmStat.Total)},
		{Name: "FreeMemory", Type: "gauge", Value: float64(vmStat.Free)},
	}

	for i, cpuUtil := range cpuPercentages {
		metrics = append(metrics, types.Metric{
			Name:  "CPUutilization" + string(i+1),
			Type:  "gauge",
			Value: cpuUtil,
		})
	}

	return metrics
}
