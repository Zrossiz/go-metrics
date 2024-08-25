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
)

func TestMetricGauge(t *testing.T) {

	store := memstorage.NewMemStorage()

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

	Metric(rr, req, store)

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

func TestMetricCounter(t *testing.T) {

	store := memstorage.NewMemStorage()

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

	Metric(rr, req, store)

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
