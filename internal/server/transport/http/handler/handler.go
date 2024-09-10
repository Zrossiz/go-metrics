package handler

import (
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/internal/server/service"
	"go.uber.org/zap"
)

type MetricHandler struct {
	service *service.MetricService
	logger  *zap.Logger
}

// Объявить интерфейс для работы со всеми функциями из сервисного слоя
type MetricHandlerer interface {
	CreateCounter(body dto.PostMetricDto) error
	Get() (models.Metric, error)
	GreateGauge(body dto.PostMetricDto) error
	GetHTML()
}

func New(s *service.MetricService, logger *zap.Logger) MetricHandler {
	return MetricHandler{
		service: s,
		logger:  logger,
	}
}

func (m *MetricHandler) Test() {

}
