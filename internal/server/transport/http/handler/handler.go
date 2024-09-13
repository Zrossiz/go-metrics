package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MetricHandler struct {
	service *service.MetricService
	logger  *zap.Logger
}

type MetricHandlerer interface {
	GetHTML(rw http.ResponseWriter, r *http.Request)
	CreateParamMetric(rw http.ResponseWriter, r *http.Request)
	CreateJSONMetric(rw http.ResponseWriter, r *http.Request)
	GetStringMetric(rw http.ResponseWriter, r *http.Request)
	GetJSONMetric(rw http.ResponseWriter, r *http.Request)
}

func New(s *service.MetricService, logger *zap.Logger) MetricHandler {
	return MetricHandler{
		service: s,
		logger:  logger,
	}
}

func (m *MetricHandler) CreateParamMetric(rw http.ResponseWriter, r *http.Request) {
	dto := dto.PostMetricDto{
		Name: chi.URLParam(r, "name"),
		Type: chi.URLParam(r, "type"),
	}

	valueMetric := chi.URLParam(r, "value")

	switch dto.Type {
	case models.GaugeType:
		float64MetricValue, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			http.Error(rw, "invalid value metric", http.StatusBadRequest)
			return
		}
		dto.Value = float64MetricValue
	case models.CounterType:
		int64MetricValue, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			http.Error(rw, "invalid value metric", http.StatusBadRequest)
			return
		}
		dto.Value = float64(int64MetricValue)
	}

	err := m.service.Create(dto)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "create metric error", http.StatusInternalServerError)
		return
	}

	metric, err := m.service.Get(dto.Name)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "get created metric error", http.StatusInternalServerError)
		return
	}

	responseString := fmt.Sprintf("Type: %s, Name: %s, Value: %v", metric.Type, metric.Name, dto.Value)

	io.WriteString(rw, responseString)
}

func (m *MetricHandler) CreateJSONMetric(rw http.ResponseWriter, r *http.Request) {
	var body dto.PostMetricDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Print(body)
	fmt.Println("")
	err = m.service.Create(body)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "create metric error", http.StatusInternalServerError)
		return
	}

	metric, err := m.service.Get(body.Name)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "get created metric error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(metric)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func (m *MetricHandler) GetStringMetric(rw http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "name")

	metric, err := m.service.GetStringValueMetric(nameMetric)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "error get metric", http.StatusInternalServerError)
		return
	}

	io.WriteString(rw, metric)
}

func (m *MetricHandler) GetJSONMetric(rw http.ResponseWriter, r *http.Request) {
	var body dto.GetMetricDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	metric, err := m.service.Get(body.Name)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "get metric error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(metric)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func (m *MetricHandler) GetHTML(rw http.ResponseWriter, r *http.Request) {
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
				<td>
					{{if eq .Type "gauge"}}
						{{.Value}}
					{{else}}
						{{.Delta}}
					{{end}}
				</td>
			</tr>
			{{end}}
		</table>
	</body>
	</html>
`

	t, err := template.New("metrics").Parse(tmpl)
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	metrics, err := m.service.GetAll()
	if err != nil {
		m.logger.Error("internal error", zap.Error(err))
		http.Error(rw, "internal server error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(rw, metrics); err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (m *MetricHandler) PingDB(rw http.ResponseWriter, r *http.Request) {
	err := m.service.PingDB()
	if err != nil {
		http.Error(rw, "connection not available", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
