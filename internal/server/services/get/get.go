package get

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func JSONMetric(rw http.ResponseWriter, r *http.Request, store memstorage.MemStorage) {
	var body dto.GetMetricDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	metric := store.GetMetric(body.ID)
	if metric == nil {
		http.Error(rw, "metric not found", http.StatusNotFound)
	}

	response, err := json.Marshal(metric)
	if err != nil {
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func Metric(rw http.ResponseWriter, r *http.Request, store memstorage.MemStorage) {
	nameMetric := chi.URLParam(r, "name")
	metric := store.GetMetric(nameMetric)

	if metric.Name == "" {
		http.Error(rw, "metric not found", 404)
		return
	}

	io.WriteString(rw, fmt.Sprintf("%v", metric.Value))
}

func HTMLPageMetric(rw http.ResponseWriter, r *http.Request, store memstorage.MemStorage) {
	metrics := store.Metrics
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Метрики</title>
		</head>
		<body>
			<h1>Список метрик</h1>
			<table border="1">
				<tr>
					<th>Имя метрики</th>
					<th>Тип метрики</th>
					<th>Значение</th>
				</tr>
				{{range .}}
				<tr>
					<td>{{.Name}}</td>
					<td>{{.Type}}</td>
					<td>{{.Value}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>
	`

	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(rw, metrics); err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}
