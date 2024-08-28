package app

import (
	"fmt"
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/logger"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger/request"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func StartServer() error {
	config.FlagParse()

	r := chi.NewRouter()

	store := memstorage.NewMemStorage()

	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		return err
	}

	zLogger := logger.Log

	r.Use(func(next http.Handler) http.Handler {
		return request.WithLogs(next)
	})

	r.Use(func(next http.Handler) http.Handler {
		return gzip.DecompressMiddleware(next)
	})

	r.Use(func(next http.Handler) http.Handler {
		return gzip.CompressMiddleware(next)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		get.HTMLPageMetric(w, r, *store)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
			update.Metric(w, r, store)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			update.JSONMetric(w, r, store)
		})
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			get.Metric(w, r, *store)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			get.JSONMetric(w, r, *store)
		})
	})

	zLogger.Info("Starting server",
		zap.String("address", config.RunAddr),
	)
	if err := http.ListenAndServe(config.RunAddr, r); err != nil {
		fmt.Println(err)
	}

	return nil
}
