package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Моковые структуры для тестирования
type MockHandler struct{}

func (m *MockHandler) GetHTML(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("HTML Response"))
}

func (m *MockHandler) CreateParamMetric(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Param Metric Created"))
}

func (m *MockHandler) CreateJSONMetric(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("JSON Metric Created"))
}

func (m *MockHandler) GetStringMetric(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("String Metric"))
}

func (m *MockHandler) GetJSONMetric(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(`{"metric": "value"}`))
}

func (m *MockHandler) PingDB(rw http.ResponseWriter, _ *http.Request) {
	rw.Write([]byte("Pong"))
}

func (m *MockHandler) CreateBatchJSONMetrics(rw http.ResponseWriter, _ *http.Request) {
	rw.Write([]byte("Batch JSON Metrics Created"))
}

// Моковый логгер
func newLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestGetHTML(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "HTML Response", rec.Body.String())
}

func TestPingDB(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Pong", rec.Body.String())
}

func TestCreateParamMetric(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodPost, "/update/param/test_metric/123", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Param Metric Created", rec.Body.String())
}

func TestGetStringMetric(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodGet, "/value/string/test_metric", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "String Metric", rec.Body.String())
}

func TestGetJSONMetric(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodPost, "/value", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"metric": "value"}`, rec.Body.String())
}

func TestCreateBatchJSONMetrics(t *testing.T) {
	handler := &MockHandler{}
	log := newLogger()

	r := New(handler, log)

	req := httptest.NewRequest(http.MethodPost, "/updates/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Batch JSON Metrics Created", rec.Body.String())
}
