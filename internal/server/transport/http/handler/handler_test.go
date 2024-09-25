package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/handler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockMetricService должен реализовывать интерфейс service.MetricServiceer
type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) Create(body dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockMetricService) Get(name string) (*models.Metric, error) {
	args := m.Called(name)
	if metric, ok := args.Get(0).(*models.Metric); ok {
		return metric, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMetricService) GetAll() (*[]models.Metric, error) {
	args := m.Called()
	if metrics, ok := args.Get(0).(*[]models.Metric); ok {
		return metrics, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMetricService) GetStringValueMetric(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockMetricService) PingDB() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMetricService) SetBatch(body []dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func TestMetricHandler_CreateParamMetric(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	mockService.On("Create", mock.Anything).Return(nil)

	mockMetric := &models.Metric{
		Name:  "TestGauge",
		Type:  models.GaugeType,
		Value: new(float64),
	}
	*mockMetric.Value = 123.45
	mockService.On("Get", "TestGauge").Return(mockMetric, nil)

	req := httptest.NewRequest("POST", "/update/gauge/TestGauge/123.45", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("name", "TestGauge")
	rctx.URLParams.Add("value", "123.45")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.CreateParamMetric(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody := "Type: gauge, Name: TestGauge, Value: 123.45"
	assert.Equal(t, expectedBody, strings.TrimSpace(rr.Body.String()))

	expectedDto := dto.PostMetricDto{
		ID:    "TestGauge",
		MType: models.GaugeType,
		Value: new(float64),
	}
	*expectedDto.Value = 123.45
	mockService.AssertCalled(t, "Create", expectedDto)
}
