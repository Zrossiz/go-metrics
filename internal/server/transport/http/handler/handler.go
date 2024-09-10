package handler

import (
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"go.uber.org/zap"
)

type MetricHandlers struct {
	service storage.Storage
	logger  *zap.Logger
}

func New(s storage.Storage, logger *zap.Logger) MetricHandlers {
	return MetricHandlers{
		service: s,
		logger:  logger,
	}
}
