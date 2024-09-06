package router

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger/request"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	"github.com/go-chi/chi/v5"
)

var ChiRouter *chi.Mux

func InitRouter() {
	ChiRouter = chi.NewRouter()

	// Применение middleware для логирования запросов
	ChiRouter.Use(func(next http.Handler) http.Handler {
		return request.WithLogs(next)
	})

	// Применение middleware для декомпрессии запросов
	ChiRouter.Use(func(next http.Handler) http.Handler {
		return gzip.DecompressMiddleware(next)
	})

	// Применение middleware для компрессии ответов
	ChiRouter.Use(func(next http.Handler) http.Handler {
		return gzip.CompressMiddleware(next)
	})

	ChiRouter.Get("/", get.HTMLPageMetric)

	ChiRouter.Get("/ping", get.Ping)

	ChiRouter.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", update.Metric)

		r.Post("/", update.JSONMetric)
	})

	ChiRouter.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", get.Metric)

		r.Post("/", get.JSONMetric)
	})
}
