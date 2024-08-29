package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/logger"
	"github.com/Zrossiz/go-metrics/internal/server/libs/parser"
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

	if config.Restore {
		_, err := parser.CollectMetricsFromFile(config.FileStoragePath, store)
		if err != nil {
			zLogger.Sugar().Fatalf("Error collect metrics %v", err)
		}
	}

	ticker := time.NewTicker(time.Duration(config.StoreInterval * int(time.Second)))
	defer ticker.Stop()
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				zLogger.Info("save metrics...")
				parser.UpdateMetrics(config.FileStoragePath, zLogger, store)
			case <-stop:
				fmt.Println("Stopping task execution")
				return
			}
		}
	}()

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
	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: r,
	}

	go func() {
		zLogger.Info("Starting server", zap.String("address", config.RunAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zLogger.Fatal("ListenAndServe():", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zLogger.Info("Shutting down server...")

	zLogger.Info("Collect metrics...")

	err := shutdownServer(store)
	if err != nil {
		zLogger.Error("Failed to save metrics on shutdown:", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zLogger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	zLogger.Info("Server exiting")
	return nil
}

func shutdownServer(store *memstorage.MemStorage) error {
	err := parser.UpdateMetrics(config.FileStoragePath, logger.Log, store)
	if err != nil {
		return err
	}

	return nil
}
