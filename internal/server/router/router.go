package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/handler"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger/request"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	"github.com/go-chi/chi/v5"
)

var ChiRouter *chi.Mux

func InitRouter() {
	ChiRouter = chi.NewRouter()

	ChiRouter.Use(func(next http.Handler) http.Handler {
		return request.WithLogs(next)
	})

	ChiRouter.Use(func(next http.Handler) http.Handler {
		return gzip.DecompressMiddleware(next)
	})

	ChiRouter.Use(func(next http.Handler) http.Handler {
		return gzip.CompressMiddleware(next)
	})

	ChiRouter.Get("/", handler.GetHTMLPageMetrics)

	ChiRouter.Get("/ping", handler.PingDB)

	ChiRouter.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", update.Metric)
		r.Post("/", update.JSONMetric)
	})

	ChiRouter.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handler.GetStringValueMetric)

		r.Post("/", handler.GetJSONMetric)
	})
}
