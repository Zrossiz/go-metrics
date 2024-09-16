package send

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/agent/constants/types"
	"github.com/Zrossiz/go-metrics/internal/agent/dto"
	"github.com/stretchr/testify/assert"
)

func createTestMetrics() []types.Metric {
	return []types.Metric{
		{
			Name:  "testMetric1",
			Type:  "gauge",
			Value: 100.0,
		},
		{
			Name:  "testMetric2",
			Type:  "counter",
			Value: 42.0,
		},
	}
}

func TestMetrics(t *testing.T) {
	metrics := createTestMetrics()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sentMetrics := Metrics(metrics, server.URL[7:])
	assert.Equal(t, len(metrics), len(*sentMetrics), "All metrics should be sent successfully")
}

func TestGzipMetrics(t *testing.T) {
	metrics := createTestMetrics()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.Path)

		contentEncoding := r.Header.Get("Content-Encoding")
		assert.Equal(t, "gzip", contentEncoding, "Content should be gzipped")

		gz, err := gzip.NewReader(r.Body)
		assert.NoError(t, err)
		defer gz.Close()

		body, err := io.ReadAll(gz)
		assert.NoError(t, err)

		var received dto.PostMetricDTO
		err = json.Unmarshal(body, &received)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sentMetrics := GzipMetrics(metrics, server.URL[7:])
	assert.Equal(t, len(metrics), len(sentMetrics), "All metrics should be sent successfully")
}

func TestBatchGzipMetrics(t *testing.T) {
	metrics := createTestMetrics()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/updates/", r.URL.Path)

		contentEncoding := r.Header.Get("Content-Encoding")
		assert.Equal(t, "gzip", contentEncoding, "Content should be gzipped")

		gz, err := gzip.NewReader(r.Body)
		assert.NoError(t, err)
		defer gz.Close()

		body, err := io.ReadAll(gz)
		assert.NoError(t, err)

		var received []types.Metric
		err = json.Unmarshal(body, &received)
		assert.NoError(t, err)

		assert.Equal(t, len(metrics), len(received), "Should receive all metrics")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	BatchGzipMetrics(metrics, server.URL[7:])
}

func TestSendWithRetry(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < maxRetries {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	var data bytes.Buffer
	request, err := getRequest("POST", server.URL, data)
	assert.NoError(t, err)

	err = sendWithRetry(request)
	assert.NoError(t, err, "Request should succeed after retries")
	assert.Equal(t, maxRetries, attempts, "Should attempt the request exactly maxRetries times")
}
