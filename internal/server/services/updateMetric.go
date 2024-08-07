package services

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func UpdateMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	if typeMetric == storage.GaugeType {
		float64MetricValue, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(rw, "parsing float value error", http.StatusBadRequest)
		}

		memstorage.SetGauge(nameMetric, float64MetricValue)
	}

	if typeMetric == storage.CounterType {
		int64MetricValue, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			http.Error(rw, "parsing int value error", http.StatusBadRequest)
		}

		memstorage.SetCounter(nameMetric, int64MetricValue)
	}

	io.WriteString(rw, fmt.Sprintf("Type: %s, Name: %s, Value: %s", typeMetric, nameMetric, valueMetric))
}
