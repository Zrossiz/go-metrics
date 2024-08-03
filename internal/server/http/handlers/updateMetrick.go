package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	CounterType = "counter"
	GaugeType   = "gauge"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	pathArr := strings.Split(r.URL.Path, "/")

	if len(pathArr) != 5 {
		http.Error(w, "invalid url format", http.StatusBadRequest)
		return
	}

	typeMetric := pathArr[2]
	if typeMetric != CounterType && typeMetric != GaugeType {
		http.Error(w, "invalid type metric", http.StatusBadRequest)
		return
	}

	nameMetric := pathArr[3]
	if len(nameMetric) == 0 {
		http.Error(w, "name metric is missing", http.StatusNotFound)
		return
	}

	valueMetric := pathArr[4]
	if typeMetric == CounterType {
		_, err := strconv.Atoi(valueMetric)
		if err != nil {
			http.Error(w, "value metric is invalid", http.StatusBadRequest)
		}
	}

	if typeMetric == GaugeType {
		_, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(w, "value metric is invalid", http.StatusBadRequest)
		}
	}

	response := fmt.Sprintf("Type: %s, Name: %s, Value: %s", typeMetric, nameMetric, valueMetric)
	fmt.Println(response)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
