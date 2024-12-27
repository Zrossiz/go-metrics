package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

type MockSecurity struct {
	mock.Mock
}

func (m *MockSecurity) DecryptedFromContext(ctx context.Context) []byte {
	args := m.Called(ctx)
	return args.Get(0).([]byte)
}

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

func (m *MockMetricService) GetAll() ([]models.Metric, error) {
	args := m.Called()
	if metrics, ok := args.Get(0).([]models.Metric); ok {
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

func TestGetHTML_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	mockService.On("GetAll").Return(make([]models.Metric, 2), nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.GetHTML(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "2", rr.Header().Get("Metics-Count"))
}

func TestCreateParamMetric(t *testing.T) {
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

func TestGetStringMetric(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()
	h := handler.New(mockService, logger)

	mockService.On("GetStringValueMetric", "TestMetric").Return("123.45", nil)

	req := httptest.NewRequest("GET", "/value/TestMetric", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("name", "TestMetric")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetStringMetric(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "123.45", rr.Body.String())

	mockService.AssertCalled(t, "GetStringValueMetric", "TestMetric")
}

func TestGetStringMetric_NotFound(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	mockService.On("GetStringValueMetric", "").Return("", nil)

	req := httptest.NewRequest(http.MethodGet, "/type/name", nil)
	rr := httptest.NewRecorder()

	h.GetStringMetric(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetJSONMetric(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	t.Run("Successful Get Metric", func(t *testing.T) {
		mockMetric := &models.Metric{
			Name:  "TestMetric",
			Type:  models.GaugeType,
			Value: floatPtr(123.45),
		}
		mockService.On("Get", "TestMetric").Return(mockMetric, nil)

		body := dto.GetMetricDto{ID: "TestMetric"}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/metric", bytes.NewReader(bodyBytes))
		rr := httptest.NewRecorder()

		h.GetJSONMetric(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		expectedResponse := dto.PostMetricDto{
			ID:    "TestMetric",
			MType: models.GaugeType,
			Value: floatPtr(123.45),
		}
		expectedResponseBytes, _ := json.Marshal(expectedResponse)
		assert.JSONEq(t, string(expectedResponseBytes), rr.Body.String())

		mockService.AssertCalled(t, "Get", "TestMetric")
	})
}

func TestCreateParamMetric_InvalidMetricType(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	req := httptest.NewRequest("POST", "/update/invalidMetricType/TestMetric/123.45", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "invalidMetricType")
	rctx.URLParams.Add("name", "TestMetric")
	rctx.URLParams.Add("value", "123.45")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.CreateParamMetric(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expectedBody := "invalid metric type"
	assert.Equal(t, expectedBody, strings.TrimSpace(rr.Body.String()))
}

func TestCreateParamMetric_InvalidValue(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	req := httptest.NewRequest("POST", "/update/gauge/TestMetric/invalidValue", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("name", "TestMetric")
	rctx.URLParams.Add("value", "invalidValue")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.CreateParamMetric(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	expectedBody := "invalid value metric"
	assert.Equal(t, expectedBody, strings.TrimSpace(rr.Body.String()))
}

func TestPingDB_Failure(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()
	h := handler.New(mockService, logger)

	mockService.On("PingDB").Return(fmt.Errorf("db not available"))

	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	h.PingDB(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "connection not available", strings.TrimSpace(rr.Body.String()))
}

func TestPingDB_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()
	h := handler.New(mockService, logger)

	mockService.On("PingDB").Return(nil)

	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	h.PingDB(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateParamMetric_ServiceError(t *testing.T) {
	mockService := new(MockMetricService)
	logger := zap.NewExample()

	h := handler.New(mockService, logger)

	mockService.On("Create", mock.Anything).Return(fmt.Errorf("service error"))

	req := httptest.NewRequest("POST", "/update/gauge/TestMetric/123.45", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("name", "TestMetric")
	rctx.URLParams.Add("value", "123.45")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.CreateParamMetric(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "create metric error")
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
