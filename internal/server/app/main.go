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
	"github.com/Zrossiz/go-metrics/internal/server/router"
	filestorage "github.com/Zrossiz/go-metrics/internal/server/storage/fileStorage"
	memstorage "github.com/Zrossiz/go-metrics/internal/server/storage/memStorage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/postgres"
	"go.uber.org/zap"
)

func StartServer() error {
	// Инициализация хранилища в памяти
	memstorage.NewMemStorage()
	store := memstorage.MemStore

	// Инициализация логгера
	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		zap.L().Error("Failed to initialize logger", zap.Error(err))
		return err
	}

	// Парсим флаги и переменный окружения
	err := config.FlagParse()
	if err != nil {
		logger.Log.Sugar().Panicf("error parse config: ", err)
	}

	// Создаем подключение к базе данных
	logger.Log.Info("connect to db...")
	err = postgres.InitConnect(config.DBConnString)
	if err != nil {
		logger.Log.Sugar().Panicf("error connect to db", err)
	}
	logger.Log.Info("db connected")

	// Восстановление метрик из файла, если включена опция Restore
	if config.Restore {
		logger.Log.Info("Restoring metrics from file", zap.String("file", config.FileStoragePath))
		_, err := filestorage.CollectMetricsFromFile(config.FileStoragePath, store)
		if err != nil {
			logger.Log.Sugar().Fatalf("Failed to collect metrics from file: %v", err)
		}
	}

	// Настройка тикера для сохранения метрик
	ticker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer func() {
		logger.Log.Info("Stopping ticker")
		ticker.Stop()
	}()
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Log.Info("Saving metrics to file", zap.String("file", config.FileStoragePath))
				if err := filestorage.UpdateMetrics(config.FileStoragePath, store); err != nil {
					logger.Log.Error("Failed to save metrics to file", zap.Error(err))
				}
			case <-stop:
				logger.Log.Info("Stopping task execution")
			}
		}
	}()

	router.InitRouter()

	// Инициализация HTTP сервера
	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: router.ChiRouter,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		logger.Log.Info("Starting server", zap.String("address", config.RunAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Обработка сигнала завершения работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	// Сохранение метрик при завершении работы сервера
	if err := shutdownServer(store); err != nil {
		logger.Log.Error("Failed to save metrics on shutdown", zap.Error(err))
	}

	// Завершение работы сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited")
	return nil
}

func shutdownServer(store *memstorage.MemStorage) error {
	logger.Log.Info("Saving metrics to file during shutdown", zap.String("file", config.FileStoragePath))
	err := filestorage.UpdateMetrics(config.FileStoragePath, store)
	if err != nil {
		logger.Log.Error("Failed to save metrics to file", zap.Error(err))
		return err
	}

	return nil
}
