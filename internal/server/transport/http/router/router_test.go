package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Мок-структура для тестирования обработчиков
type mockMetricRouter struct{}

func (m *mockMetricRouter) GetHTML(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HTML response"))
}

func (m *mockMetricRouter) CreateParamMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created param metric"))
}

func (m *mockMetricRouter) CreateJSONMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created JSON metric"))
}

func (m *mockMetricRouter) GetStringMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric value"))
}

func (m *mockMetricRouter) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"metric": "value"}`))
}

func (m *mockMetricRouter) PingDB(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DB is reachable"))
}

func (m *mockMetricRouter) CreateBatchJSONMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Batch JSON metrics created"))
}

func TestMetricRouter(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	rr := New(&mockMetricRouter{}, logger)

	tests := []struct {
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			method:       http.MethodGet,
			url:          "/",
			expectedCode: http.StatusOK,
			expectedBody: "HTML response",
		},
		{
			method:       http.MethodGet,
			url:          "/ping",
			expectedCode: http.StatusOK,
			expectedBody: "DB is reachable",
		},
		{
			method:       http.MethodPost,
			url:          "/update/gauge/metric/10",
			expectedCode: http.StatusOK,
			expectedBody: "Created param metric",
		},
		{
			method:       http.MethodPost,
			url:          "/update",
			expectedCode: http.StatusOK,
			expectedBody: "Created JSON metric",
		},
		{
			method:       http.MethodPost,
			url:          "/updates/",
			expectedCode: http.StatusOK,
			expectedBody: "Batch JSON metrics created",
		},
		{
			method:       http.MethodGet,
			url:          "/value/gauge/metric",
			expectedCode: http.StatusOK,
			expectedBody: "Metric value",
		},
		{
			method:       http.MethodPost,
			url:          "/value",
			expectedCode: http.StatusOK,
			expectedBody: `{"metric": "value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)

			rrw := httptest.NewRecorder()
			rr.ServeHTTP(rrw, req)

			assert.Equal(t, tt.expectedCode, rrw.Code)
			assert.Contains(t, rrw.Body.String(), tt.expectedBody)
		})
	}
}
