package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MetricRouter interface {
	GetHTML(rw http.ResponseWriter, r *http.Request)
	CreateParamMetric(rw http.ResponseWriter, r *http.Request)
	CreateJSONMetric(rw http.ResponseWriter, r *http.Request)
	GetStringMetric(rw http.ResponseWriter, r *http.Request)
	GetJSONMetric(rw http.ResponseWriter, r *http.Request)
	PingDB(rw http.ResponseWriter, _ *http.Request)
	CreateBatchJSONMetrics(rw http.ResponseWriter, _ *http.Request)
}

func New(handl MetricRouter, log *zap.Logger) http.Handler {

	r := chi.NewRouter()

	r.Use(logger.WithLogs(log))
	r.Use(gzip.DecompressMiddleware)
	r.Use(gzip.CompressMiddleware)

	r.Get("/", handl.GetHTML)

	r.Get("/ping", handl.PingDB)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handl.CreateParamMetric)
		r.Post("/", handl.CreateJSONMetric)
	})

	r.Post("/updates/", handl.CreateBatchJSONMetrics)

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handl.GetStringMetric)

		r.Post("/", handl.GetJSONMetric)
	})

	return r
}
