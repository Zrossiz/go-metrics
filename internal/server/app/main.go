package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/logger"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/gzip"
	"github.com/Zrossiz/go-metrics/internal/server/middleware/logger/request"
	"github.com/Zrossiz/go-metrics/internal/server/services/get"
	"github.com/Zrossiz/go-metrics/internal/server/services/update"
	filestorage "github.com/Zrossiz/go-metrics/internal/server/storage/fileStorage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func StartServer() error {
	// Разбор флагов конфигурации
	config.FlagParse()

	r := chi.NewRouter()

	// Инициализация хранилища в памяти
	store := memstorage.NewMemStorage()

	// Инициализация логгера
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		zap.L().Error("Failed to initialize logger", zap.Error(err))
		return err
	}

	zLogger := logger.Log

	// Восстановление метрик из файла, если включена опция Restore
	if config.Restore {
		zLogger.Info("Restoring metrics from file", zap.String("file", config.FileStoragePath))
		_, err := filestorage.CollectMetricsFromFile(config.FileStoragePath, store)
		if err != nil {
			zLogger.Sugar().Fatalf("Failed to collect metrics from file: %v", err)
		}
	}

	// Настройка тикера для сохранения метрик
	ticker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer func() {
		zLogger.Info("Stopping ticker")
		ticker.Stop()
	}()
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				zLogger.Info("Saving metrics to file", zap.String("file", config.FileStoragePath))
				if err := filestorage.UpdateMetrics(config.FileStoragePath, zLogger, store); err != nil {
					zLogger.Error("Failed to save metrics to file", zap.Error(err))
				}
			case <-stop:
				zLogger.Info("Stopping task execution")
			}
		}
	}()

	// Применение middleware для логирования запросов
	r.Use(func(next http.Handler) http.Handler {
		return request.WithLogs(next)
	})

	// Применение middleware для декомпрессии запросов
	r.Use(func(next http.Handler) http.Handler {
		return gzip.DecompressMiddleware(next)
	})

	// Применение middleware для компрессии ответов
	r.Use(func(next http.Handler) http.Handler {
		return gzip.CompressMiddleware(next)
	})

	// Определение маршрутов
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

	// Инициализация HTTP сервера
	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: r,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		zLogger.Info("Starting server", zap.String("address", config.RunAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Обработка сигнала завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zLogger.Info("Shutting down server...")

	// Сохранение метрик при завершении работы сервера
	if err := shutdownServer(store); err != nil {
		zLogger.Error("Failed to save metrics on shutdown", zap.Error(err))
	}

	// Завершение работы сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zLogger.Info("Server exited")
	return nil
}

func shutdownServer(store *memstorage.MemStorage) error {
	zLogger := logger.Log

	zLogger.Info("Saving metrics to file during shutdown", zap.String("file", config.FileStoragePath))
	err := filestorage.UpdateMetrics(config.FileStoragePath, zLogger, store)
	if err != nil {
		zLogger.Error("Failed to save metrics to file", zap.Error(err))
		return err
	}

	return nil
}
