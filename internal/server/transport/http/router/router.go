package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func New(s *service.MetricService, log *zap.Logger) http.Handler {
	handl := handler.New(s, log)

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

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handl.GetStringMetric)

		r.Post("/", handl.GetJSONMetric)
	})

	return r
}
