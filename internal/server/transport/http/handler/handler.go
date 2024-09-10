package handler

import (
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"go.uber.org/zap"
)

type MetricHandlers struct {
	service *service.MetricService
	logger  *zap.Logger
}

func New(s *service.MetricService, logger *zap.Logger) MetricHandlers {
	return MetricHandlers{
		service: s,
		logger:  logger,
	}
}

func (m *MetricHandlers) Test() {

}
