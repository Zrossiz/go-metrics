package get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
)

func TestMetric(t *testing.T) {
	memstorage.NewMemStorage()
	expectedValue := 42.42

	memstorage.MemStore.Metrics = []storage.Metric{
		{Name: "testMetric", Type: storage.GaugeType, Value: expectedValue},
	}

	mockMetric := dto.GetMetricDto{
		ID:    "testMetric",
		MType: "gauge",
	}

	jsonMetric, err := json.Marshal(mockMetric)
	if err != nil {
		fmt.Println("json metric parsing error")
	}

	req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(jsonMetric))

	rr := httptest.NewRecorder()

	JSONMetric(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}

	var body dto.MetricDTO

	err = json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Errorf("parsing response error")
	}

	if *body.Value != expectedValue {
		t.Errorf("expected value metric %v, got %v", expectedValue, body.Value)
	}
}
