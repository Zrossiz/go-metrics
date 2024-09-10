package app

import (
	"time"

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
	}

	store := storage.New(dbConn, cfg.FileStoragePath, time.Duration(cfg.StoreInterval), log.ZapLogger)

	serv := service.New(store, log.ZapLogger)

	r := router.New(serv, log.ZapLogger)
}
