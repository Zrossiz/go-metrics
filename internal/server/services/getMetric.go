package services

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
)

func GetMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "type")
	nameMetric := chi.URLParam(r, "name")
	metric := memstorage.GetMetric(nameMetric)

	if metric.Name == "" {
		http.Error(rw, "metric not found", 404)
		return
	}

	io.WriteString(rw, fmt.Sprintf("Type: %s, Name: %s, Value: %s", typeMetric, nameMetric, metric.Value))
}

func GetHTMLPageMetric(rw http.ResponseWriter, r *http.Request) {
	metrics := memstorage.Metrics
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
