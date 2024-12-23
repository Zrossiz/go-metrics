// Package app is responsible for configuring and starting the server
// as well as handling graceful shutdowns and monitoring
package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // For performance profiling
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/security"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/handler"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/router"
	"github.com/Zrossiz/go-metrics/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// StartServer initializes and starts the HTTP server.
// It handles the following steps:
//   - Parses configuraiton.
//   - Initializes the logger.
//   - Configures database connection (if provided).
//   - Sets up storage, service and transport layers.
//   - Starts the HTTP server and the `pprof` monitoring server.
//   - Handles graceful shutdown on receiving system signals (SIGINT or SIGTERM).
func StartServer() {
	// Parse configuration
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("get config error", zap.Error(err))
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		fmt.Println("init logger error", zap.Error(err))
	}

	// Configure database connection (if DSN is provided)
	var dbConn *pgxpool.Pool
	if len(cfg.DBDSN) > 0 {
		dbConn, err = dbstorage.GetConnect(cfg.DBDSN, log.ZapLogger)
		if err != nil {
			log.ZapLogger.Fatal("error connect to db", zap.Error(err))
		}
	}

	// Configure crypto key
	cryptoKey, err := security.GetPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		log.ZapLogger.Fatal("get crypto key error")
	}
	cfg.PrivateKey = cryptoKey
	config.AppConfig.PrivateKey = cryptoKey

	// Initialize the storage layer
	store := storage.New(dbConn, cfg, log.ZapLogger)

	// Initialize the service layer (business logic)
	serv := service.New(store)

	// Initialize the transport layer (HTTP handlers)
	handl := handler.New(serv, log.ZapLogger)
	r := router.New(&handl, log.ZapLogger)

	// Configure the HTTP server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	// Start the HTTP server in a goroutine
	go func() {
		log.ZapLogger.Info("Starting server", zap.String("address", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ZapLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Start the `pprof` monitoring server in a separete goroutine
	go func() {
		log.ZapLogger.Info("Starting pprof server on localhost:6060")
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.ZapLogger.Error("Failed to start pprof server", zap.Error(err))
		}
	}()

	// Wait for shutdown signal (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.ZapLogger.Info("Shutting down server...")

	// Handle greaceful shutdown and save state if necessary
	if err := shutdownServer(store, log.ZapLogger, *cfg); err != nil {
		log.ZapLogger.Error("Failed to save metrics on shutdown", zap.Error(err))
	}

	// Gracefully stop the HTTP server
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
