package update

import (
	"encoding/json"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
)

func Metric(rw http.ResponseWriter, r *http.Request, store *memstorage.MemStorage) {
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

	response, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
