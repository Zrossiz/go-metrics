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

	ChiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		get.HTMLPageMetric(w, r)
	})

	ChiRouter.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
			update.Metric(w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			update.JSONMetric(w, r)
		})
	})

	ChiRouter.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			get.Metric(w, r)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			get.JSONMetric(w, r)
		})
	})
}
