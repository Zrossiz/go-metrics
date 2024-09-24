package collector

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/constants"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type memMock struct {
	mock.Mock
}

func (m *memMock) VirtualMemory() (*mem.VirtualMemoryStat, error) {
	args := m.Called()
	return args.Get(0).(*mem.VirtualMemoryStat), args.Error(1)
}

// Mock for cpu.Percent
type cpuMock struct {
	mock.Mock
}

func (m *cpuMock) Percent(interval float64, percpu bool) ([]float64, error) {
	args := m.Called(interval, percpu)
	return args.Get(0).([]float64), args.Error(1)
}

func TestGetMetrics(t *testing.T) {
	// Создаем mock для памяти
	memMock := new(memMock)
	memMock.On("VirtualMemory").Return(&mem.VirtualMemoryStat{
		Total: 4096,
		Free:  1024,
	}, nil)

	// Создаем mock для CPU
	cpuMock := new(cpuMock)
	cpuMock.On("Percent", 0.0, true).Return([]float64{10.5, 15.2}, nil)

	// Счетчик для теста
	var counter int64 = 0

	// Вызов GetMetrics
	metrics := GetMetrics(&counter)

	// Проверяем результат
	assert.NotNil(t, metrics)
	assert.Greater(t, len(metrics), 0)

	// Проверяем, что некоторые ключевые метрики присутствуют
	assert.Equal(t, constants.Counter, metrics[0].Type)
	assert.Equal(t, "PollCount", metrics[0].Name)
	assert.Equal(t, float64(1), metrics[0].Value) // Счетчик был увеличен

	// Проверяем метрики памяти
	foundTotalMemory := false
	foundFreeMemory := false
	for _, metric := range metrics {
		if metric.Name == "TotalMemory" {
			foundTotalMemory = true
			assert.Equal(t, float64(4096), metric.Value)
		}
		if metric.Name == "FreeMemory" {
			foundFreeMemory = true
			assert.Equal(t, float64(1024), metric.Value)
		}
	}
	assert.True(t, foundTotalMemory)
	assert.True(t, foundFreeMemory)

	// Проверяем метрики CPU
	foundCPU1 := false
	foundCPU2 := false
	for _, metric := range metrics {
		if metric.Name == "CPUutilization1" {
			foundCPU1 = true
			assert.Equal(t, float64(10.5), metric.Value)
		}
		if metric.Name == "CPUutilization2" {
			foundCPU2 = true
			assert.Equal(t, float64(15.2), metric.Value)
		}
	}
	assert.True(t, foundCPU1)
	assert.True(t, foundCPU2)
}
