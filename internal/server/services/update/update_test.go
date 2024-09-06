package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/dto"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func TestMetricCounter(t *testing.T) {
	memstorage.NewMemStorage()
	store := memstorage.MemStore

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		Metric(w, r)
	})

	req := httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/42", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	metric := store.GetMetric("testCounter")

	if metric == nil {
		t.Errorf("expected metric to be created, got nil")
		return
	}

	if value, ok := metric.Value.(int64); !ok || value != 42 {
		t.Errorf("expected Counter value 42, got %v", metric.Value)
	}
}

func TestMetricGauge(t *testing.T) {
	memstorage.NewMemStorage()
	store := memstorage.MemStore

	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		Metric(w, r)
	})

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/testGauge/42.42", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	metric := store.GetMetric("testGauge")

	if metric == nil {
		t.Errorf("expected metric to be created, got nil")
		return
	}

	if value, ok := metric.Value.(float64); !ok || value != 42.42 {
		t.Errorf("expected Counter value 42.42, got %v", metric.Value)
	}
}

func TestJSONMetricGauge(t *testing.T) {
	memstorage.NewMemStorage()
	store := memstorage.MemStore

	mockMetricValue := 42.42

	mockMetric := dto.MetricDTO{
		ID:    "testGauge",
		MType: "gauge",
		Value: &mockMetricValue,
	}

	jsonMetric, err := json.Marshal(mockMetric)
	if err != nil {
		fmt.Println("json metric parsing error")
	}

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonMetric))

	rr := httptest.NewRecorder()

	JSONMetric(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	createdMetric := store.GetMetric("testGauge")

	if createdMetric == nil {
		t.Errorf("expected metric, got nil")
		return
	}

	if createdMetric.Value != 42.42 {
		t.Errorf("expected Gauge value 42.42, got %v", createdMetric.Value)
	}
}

func TestJSONMetricCounter(t *testing.T) {
	memstorage.NewMemStorage()
	store := memstorage.MemStore

	var mockMetricValue int64 = 42

	mockMetric := dto.MetricDTO{
		ID:    "testCounter",
		MType: "counter",
		Delta: &mockMetricValue,
	}

	jsonMetric, err := json.Marshal(mockMetric)
	if err != nil {
		fmt.Println("json metric parsing error")
	}

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonMetric))

	rr := httptest.NewRecorder()

	JSONMetric(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	createdMetric := store.GetMetric("testCounter")

	if createdMetric == nil {
		t.Errorf("expected metric, got nil")
		return
	}

	if createdMetric.Value.(int64) != 42 {
		t.Errorf("expected Counter value 42, got %v", createdMetric.Value)
	}
}
