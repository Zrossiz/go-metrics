package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/handler"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/router"
	"github.com/Zrossiz/go-metrics/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func StartServer() {
	cfg, err := config.GetConfig()
	if err != nil {
		zap.S().Fatalf("get config error", zap.Error(err))
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		zap.S().Fatalf("init logger error", zap.Error(err))
	}

	var dbConn *pgxpool.Pool
	if len(cfg.DBDSN) > 0 {
		dbConn, err = dbstorage.GetConnect(cfg.DBDSN, log.ZapLogger)
		if err != nil {
			log.ZapLogger.Fatal("error connect to db", zap.Error(err))
		}
	}

	store := storage.New(dbConn, cfg, log.ZapLogger)
	serv := service.New(store, log.ZapLogger, dbConn)
	handl := handler.New(serv, log.ZapLogger)
	r := router.New(&handl, log.ZapLogger)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	go func() {
		log.ZapLogger.Info("Starting server", zap.String("address", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ZapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.ZapLogger.Info("Shutting down server...")

	if err := shutdownServer(store, log.ZapLogger, *cfg); err != nil {
		log.ZapLogger.Error("Failed to save metrics on shutdown", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.ZapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.ZapLogger.Info("Server exited")
}

func shutdownServer(store storage.Storage, log *zap.Logger, cfg config.Config) error {
	log.Info("Saving metrics to file during shutdown", zap.String("file", cfg.FileStoragePath))
	err := store.Close()

	if err != nil {
		log.Error("Failed to close connection", zap.Error(err))
		return err
	}

	return nil
}
