package app

import (
	"net/http"

	"github.com/Zrossiz/go-metrics/internal/server/config"
	"github.com/Zrossiz/go-metrics/internal/server/libs/logger"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"github.com/Zrossiz/go-metrics/internal/server/storage/dbstorage"
	"github.com/Zrossiz/go-metrics/internal/server/transport/http/router"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	var dbConn *gorm.DB
	if len(cfg.DBDSN) > 0 {
		dbConn, err = dbstorage.GetConnect(cfg.DBDSN)
		if err != nil {
			log.ZapLogger.Fatal("error connect to db", zap.Error(err))
		}

		err = dbstorage.MigrateSQL(dbConn)
		if err != nil {
			log.ZapLogger.Fatal("migrate error", zap.Error(err))
		}
	}

	store := storage.New(dbConn, cfg.FileStoragePath)

	if cfg.Restore {
		log.ZapLogger.Info("start collect metrics from storage...")
		err := store.Load(cfg.FileStoragePath)
		if err != nil {
			log.ZapLogger.Fatal("error collect metric", zap.Error(err))
		}
		log.ZapLogger.Info("metric collected")
	}

	serv := service.New(store, log.ZapLogger)

	r := router.New(serv, log.ZapLogger)

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
}

// ticker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
// defer func() {
// 	logger.Log.Info("Stopping ticker")
// 	ticker.Stop()
// }()
// stop := make(chan bool)

// go func() {
// 	for {
// 		select {
// 		case <-ticker.C:
// 			logger.Log.Info("Saving metrics to file", zap.String("file", config.FileStoragePath))
// 			if err := filestorage.UpdateMetrics(config.FileStoragePath, store); err != nil {
// 				logger.Log.Error("Failed to save metrics to file", zap.Error(err))
// 			}
// 		case <-stop:
// 			logger.Log.Info("Stopping task execution")
// 		}
// 	}
// }()
