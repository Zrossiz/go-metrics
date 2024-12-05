package service

import (
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorager is a mock implementation of the Storager interface
type MockStorager struct {
	mock.Mock
}

func (m *MockStorager) SetGauge(body dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockStorager) SetCounter(body dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockStorager) SetBatch(body []dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockStorager) Get(name string) (*models.Metric, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Metric), args.Error(1)
}

func (m *MockStorager) GetAll() ([]models.Metric, error) {
	args := m.Called()
	return args.Get(0).([]models.Metric), args.Error(1)
}

func (m *MockStorager) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestMetricService_CreateGauge(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	gaugeMetric := dto.PostMetricDto{
		ID:    "gauge_metric",
		MType: models.GaugeType,
		Value: float64Ptr(42.5),
	}
	mockStorage.On("SetGauge", gaugeMetric).Return(nil)

	err := service.Create(gaugeMetric)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_CreateCounter(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	counterMetric := dto.PostMetricDto{
		ID:    "counter_metric",
		MType: models.CounterType,
		Delta: int64Ptr(10),
	}
	mockStorage.On("SetCounter", counterMetric).Return(nil)

	err := service.Create(counterMetric)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_CreateMissingDelta(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	counterMetric := dto.PostMetricDto{
		ID:    "counter_metric",
		MType: models.CounterType,
		Delta: nil,
	}

	err := service.Create(counterMetric)
	assert.EqualError(t, err, "delta not found")
}

func TestMetricService_CreateMissingValue(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	gaugeMetric := dto.PostMetricDto{
		ID:    "gauge_metric",
		MType: models.GaugeType,
		Value: nil,
	}

	err := service.Create(gaugeMetric)
	assert.EqualError(t, err, "value not found")
}

func TestMetricService_SetBatch(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	metrics := []dto.PostMetricDto{
		{ID: "metric1", MType: models.GaugeType, Value: float64Ptr(12.34)},
		{ID: "metric2", MType: models.CounterType, Delta: int64Ptr(5)},
	}
	expectedBatch := []dto.PostMetricDto{
		{ID: "metric1", MType: models.GaugeType, Value: float64Ptr(12.34)},
		{ID: "metric2", MType: models.CounterType, Delta: int64Ptr(5)},
	}
	mockStorage.On("SetBatch", expectedBatch).Return(nil)

	err := service.SetBatch(metrics)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_Get(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	expectedMetric := &models.Metric{
		Name:  "test_metric",
		Type:  models.GaugeType,
		Value: float64Ptr(99.9),
	}
	mockStorage.On("Get", "test_metric").Return(expectedMetric, nil)

	result, err := service.Get("test_metric")
	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_GetStringValueMetric_Gauge(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	expectedMetric := &models.Metric{
		Name:  "gauge_metric",
		Type:  models.GaugeType,
		Value: float64Ptr(42.5),
	}
	mockStorage.On("Get", "gauge_metric").Return(expectedMetric, nil)

	value, err := service.GetStringValueMetric("gauge_metric")
	assert.NoError(t, err)
	assert.Equal(t, "42.5", value)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_GetStringValueMetric_Counter(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	expectedMetric := &models.Metric{
		Name:  "counter_metric",
		Type:  models.CounterType,
		Delta: int64Ptr(100),
	}
	mockStorage.On("Get", "counter_metric").Return(expectedMetric, nil)

	value, err := service.GetStringValueMetric("counter_metric")
	assert.NoError(t, err)
	assert.Equal(t, "100", value)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_GetAll(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	expectedMetrics := []models.Metric{
		{Name: "metric1", Type: models.GaugeType, Value: float64Ptr(1.23)},
		{Name: "metric2", Type: models.CounterType, Delta: int64Ptr(10)},
	}
	mockStorage.On("GetAll").Return(expectedMetrics, nil)

	result, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
	mockStorage.AssertExpectations(t)
}

func TestMetricService_PingDB(t *testing.T) {
	mockStorage := new(MockStorager)
	service := New(mockStorage)

	mockStorage.On("Ping").Return(nil)

	err := service.PingDB()
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

// Helper functions to create pointers for float64 and int64 values
func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
