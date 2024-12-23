package send

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/stretchr/testify/assert"
)

func TestComputeHash(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		cfgKey   string
		expected string
	}{
		{
			name:     "Valid data and cfg key",
			data:     "test data",
			cfgKey:   "test_key",
			expected: "28d2d31c85bb54f7e48bfa3f5cf407a007d3957f7788e8c6d5ed38ab080fc308",
		},
		{
			name:     "Empty cfgKey",
			data:     "test data",
			cfgKey:   "",
			expected: "",
		},
		{
			name:     "Empty data and empty cfgKey",
			data:     "",
			cfgKey:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dataBuffer bytes.Buffer
			dataBuffer.WriteString(tt.data)

			got := computeHash(dataBuffer, tt.cfgKey)
			if got != tt.expected {
				t.Errorf("expected hash %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestGetBytesMetricDTO(t *testing.T) {
	tests := []struct {
		name          string
		metric        types.Metric
		expected      []byte
		expectedError bool
	}{
		{
			name: "Metric with float64 value",
			metric: types.Metric{
				Name:  "metric_float",
				Type:  "gauge",
				Value: 123.45,
			},
			expected:      []byte(`{"id":"metric_float","type":"gauge","value":123.45}`),
			expectedError: false,
		},
		{
			name: "Metric with int64 value",
			metric: types.Metric{
				Name:  "metric_int",
				Type:  "counter",
				Value: int64(100),
			},
			expected:      []byte(`{"id":"metric_int","type":"counter","delta":100}`),
			expectedError: false,
		},
		{
			name: "Metric with nil value",
			metric: types.Metric{
				Name:  "metric_nil",
				Type:  "counter",
				Value: nil,
			},
			expected:      []byte(`{"id":"metric_nil","type":"counter"}`),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBytesMetricDTO(tt.metric)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, string(tt.expected), string(got))
			}
		})
	}
}

func TestGzipMetrics(t *testing.T) {
	var expectedCounterValue int64 = 42

	metrics := []types.Metric{
		{Type: "gauge", Name: "metric1", Value: 1.23},
		{Type: "counter", Name: "metric2", Value: expectedCounterValue},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sendedMetrics := GzipMetrics(metrics, server.Listener.Addr().String(), "key")

	if len(sendedMetrics) != len(metrics) {
		fmt.Print(sendedMetrics[0].Name)
		t.Errorf("Expected %d metrics to be sent, but got %d", len(metrics), len(sendedMetrics))
	}

	for i, metric := range sendedMetrics {
		if metric != metrics[i] {
			t.Errorf("Expected metric %v, but got %v", metrics[i], metric)
		}
	}
}

func TestBatchGzipMetrics(t *testing.T) {
	var expectedCounterValue int64 = 42

	metrics := []types.Metric{
		{Type: "gauge", Name: "metric1", Value: 1.23},
		{Type: "counter", Name: "metric2", Value: expectedCounterValue},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, t *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	BatchGzipMetrics(metrics, server.Listener.Addr().String(), "key")
	t.Log("BatchGzipMetrics executed successfully")
}

func TestGetRequest(t *testing.T) {
	tests := []struct {
		name          string
		reqURL        string
		data          bytes.Buffer
		hash          string
		method        string
		expectedError bool
	}{
		{
			name:          "Valid data",
			reqURL:        "http://localhost:8080/",
			data:          *bytes.NewBuffer([]byte(`{"id": 1}`)),
			hash:          "hash",
			method:        "POST",
			expectedError: false,
		},
		{
			name:          "Invalid method",
			reqURL:        "http://localhost:8080/",
			data:          *bytes.NewBuffer([]byte(`{"id": 1}`)),
			hash:          "hash",
			method:        "Invalid",
			expectedError: true,
		},
		{
			name:          "Empty data",
			reqURL:        "http://localhost:8080/",
			data:          *bytes.NewBuffer([]byte(``)),
			hash:          "",
			method:        "POST",
			expectedError: true,
		},
		{
			name:          "Empty URL",
			reqURL:        "",
			data:          *bytes.NewBuffer([]byte(`{"id": 1}`)),
			hash:          "hash",
			method:        "POST",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := getRequest(tt.method, tt.reqURL, tt.data, tt.hash)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
			}
		})
	}
}
