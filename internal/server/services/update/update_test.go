package update

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func TestUpdateMetricGauge(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/42.42", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", storage.GaugeType)
	rctx.URLParams.Add("name", "testGauge")
	rctx.URLParams.Add("value", "42.42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	UpdateMetric(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	if metric := memstorage.GetMetric("testGauge"); metric.Value != 42.42 {
		t.Errorf("expected Gauge value 42.42, got %v", metric.Value)
	}
}

func TestUpdateMetricCounter(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/42", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", storage.CounterType)
	rctx.URLParams.Add("name", "testCounter")
	rctx.URLParams.Add("value", "42")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	UpdateMetric(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	metric := memstorage.GetMetric("testCounter")

	if value, ok := metric.Value.(int64); !ok || value != 42 {
		t.Errorf("expected Counter value 42, got %v", metric.Value)
	}
}
