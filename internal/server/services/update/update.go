package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/libs/logger"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func JSONMetric(rw http.ResponseWriter, r *http.Request, store *memstorage.MemStorage) {
	var body dto.MetricDTO

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updatedMetric *storage.Metric

	switch body.MType {
	case storage.GaugeType:
		if body.Value == nil {
			http.Error(rw, "missing value for gauge", http.StatusBadRequest)
			return
		}
		updatedMetric = store.SetGauge(body.ID, *body.Value)
	case storage.CounterType:
		if body.Delta == nil {
			http.Error(rw, "missing delta for counter", http.StatusBadRequest)
			return
		}
		updatedMetric = store.SetCounter(body.ID, *body.Delta)
	default:
		http.Error(rw, "unknown metric type", http.StatusBadRequest)
		return
	}

	if updatedMetric == nil {
		http.Error(rw, "invalid request", http.StatusBadRequest)
		return
	}

	responseMetric := dto.MetricDTO{
		ID:    updatedMetric.Name,
		MType: updatedMetric.Type,
	}

	if v, ok := updatedMetric.Value.(float64); ok {
		responseMetric.Value = &v
	}

	if d, ok := updatedMetric.Value.(int64); ok {
		responseMetric.Delta = &d
	}

	response, err := json.Marshal(responseMetric)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func Metric(rw http.ResponseWriter, r *http.Request, store *memstorage.MemStorage) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	logger.Log.Sugar().Infoln("TYPE: ", typeMetric)
	logger.Log.Sugar().Infoln("NAME: ", nameMetric)

	switch typeMetric {
	case storage.GaugeType:
		float64MetricValue, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(rw, "parsing float value error", http.StatusBadRequest)
			return
		}
		store.SetGauge(nameMetric, float64MetricValue)
	case storage.CounterType:
		int64MetricValue, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			http.Error(rw, "parsing int value error", http.StatusBadRequest)
			return
		}
		store.SetCounter(nameMetric, int64MetricValue)
	default:
		http.Error(rw, "unknown metric type", http.StatusBadRequest)
		return
	}

	io.WriteString(rw, fmt.Sprintf("Type: %s, Name: %s, Value: %s", typeMetric, nameMetric, valueMetric))
}
