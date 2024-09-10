package service

import (
	"github.com/Zrossiz/go-metrics/internal/server/storage"
	"go.uber.org/zap"
)

type MetricService struct {
	storage storage.Storage
	logger  *zap.Logger
}

func New(stor storage.Storage, logger *zap.Logger) *MetricService {
	return &MetricService{
		storage: stor,
		logger:  logger,
	}
}
