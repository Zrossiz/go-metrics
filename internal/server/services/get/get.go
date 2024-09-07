package get

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/postgres"
)

func JSONMetric(body dto.GetMetricDto) (*dto.MetricDTO, error) {
	metric := memstorage.MemStore.GetMetric(body.ID)
	if metric == nil {
		return nil, nil
	}

	responseMetric := dto.MetricDTO{
		ID:    metric.Name,
		MType: metric.Type,
	}

	if v, ok := metric.Value.(float64); ok {
		responseMetric.Value = &v
	}

	if d, ok := metric.Value.(int64); ok {
		responseMetric.Delta = &d
	}

	return &responseMetric, nil
}

func StringValueMetric(nameMetric string) string {
	metric := memstorage.MemStore.GetMetric(nameMetric)

	if metric == nil {
		return ""
	}

	return fmt.Sprintf("%v", metric.Value)
}

func HTMLPageMetric(rw http.ResponseWriter) {
	metrics := memstorage.MemStore.Metrics
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

func Ping() error {
	err := postgres.Ping(postgres.PgConn)
	if err != nil {
		return err
	}

	return nil
}
