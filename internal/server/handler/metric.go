package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
)

func GetJSONMetric(rw http.ResponseWriter, r *http.Request) {
	var body dto.GetMetricDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	metric, err := get.JSONMetric(body)
	if err != nil {

	}

	response, err := json.Marshal(metric)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func GetHTMLPageMetrics(rw http.ResponseWriter, r *http.Request) {
	get.HTMLPageMetric(rw)
}

func PingDB(rw http.ResponseWriter, r *http.Request) {
	err := get.Ping()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	return
}
