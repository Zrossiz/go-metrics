package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	"github.com/go-chi/chi/v5"
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

func GetStringValueMetric(rw http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "name")

	valueMetric := get.StringValueMetric(nameMetric)

	if valueMetric == "" {
		http.Error(rw, "metric not found", http.StatusNotFound)
		return
	}

	io.WriteString(rw, valueMetric)
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

func UpdateJSONMetric(rw http.ResponseWriter, r *http.Request) {
	var body dto.MetricDTO

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	metric := update.JSONMetric(body)

	response, err := json.Marshal(metric)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func UpdateParamsMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	valueMetric := chi.URLParam(r, "value")

	metric, err := update.ParamMetric(typeMetric, nameMetric, valueMetric)
	if err != nil {
		http.Error(rw, "internal server error", 500)
	}

	responseString := fmt.Sprintf("Type: %s, Name: %s, Value: %s", metric.MType, metric.ID, valueMetric)

	io.WriteString(rw, responseString)
}
