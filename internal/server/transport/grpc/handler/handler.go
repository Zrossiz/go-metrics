package handler

import (
	"context"
	"errors"

	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/Zrossiz/go-metrics/proto"
	"go.uber.org/zap"
)

type GrpcHandler struct {
	service MetricService
	logger  *zap.Logger
	proto.UnimplementedMetricsServer
}

type MetricService interface {
	// Create saves a single metric based on the provided DTO.
	Create(body dto.PostMetricDto) error
	// Get retrieves a metric by its name.
	Get(name string) (*models.Metric, error)
	// GetAll retrieves all available metrics.
	GetAll() ([]models.Metric, error)
	// GetStringValueMetric retrieves a metric's value as a string.
	GetStringValueMetric(name string) (string, error)
	// PingDB verifies the database connection.
	PingDB() error
	// SetBatch creates multiple metrics at once based on the provided DTO array.
	SetBatch(body []dto.PostMetricDto) error
}

func NewGrpcHandler(service MetricService, logger *zap.Logger) *GrpcHandler {
	return &GrpcHandler{
		service: service,
		logger:  logger,
	}
}

func (h *GrpcHandler) UpdateParam(ctx context.Context, req *proto.PostMetricRequest) (*proto.PostMetricResponse, error) {
	dto := dto.PostMetricDto{
		ID:    req.Id,
		MType: req.Type,
	}

	switch req.Type {
	case models.GaugeType:
		dto.Value = &req.Value
	case models.CounterType:
		dto.Delta = &req.Delta
	default:
		return nil, errors.New("invalid metric type")
	}

	if err := h.service.Create(dto); err != nil {
		h.logger.Error("failed to create or update metric", zap.Error(err))
		return nil, err
	}

	return &proto.PostMetricResponse{Success: true}, nil
}

func (h *GrpcHandler) UpdateJSON(ctx context.Context, req *proto.PostMetricRequest) (*proto.PostMetricResponse, error) {
	dto := dto.PostMetricDto{
		ID:    req.Id,
		MType: req.Type,
		Value: &req.Value,
		Delta: &req.Delta,
	}

	if err := h.service.Create(dto); err != nil {
		h.logger.Error("failed to create metric", zap.Error(err))
		return nil, err
	}

	return &proto.PostMetricResponse{Success: true}, nil
}

func (h *GrpcHandler) UpdateBatchJSON(ctx context.Context, req *proto.BatchPostMetricRequest) (*proto.BatchPostMetricResponse, error) {
	var metrics []dto.PostMetricDto

	for _, m := range req.Metrics {
		dto := dto.PostMetricDto{
			ID:    m.Id,
			MType: m.Type,
			Value: &m.Value,
			Delta: &m.Delta,
		}
		metrics = append(metrics, dto)
	}

	if err := h.service.SetBatch(metrics); err != nil {
		h.logger.Error("failed to process batch update", zap.Error(err))
		return nil, err
	}

	return &proto.BatchPostMetricResponse{Success: true}, nil
}

func (h *GrpcHandler) GetMetric(ctx context.Context, req *proto.GetMetricRequest) (*proto.GetMetricResponse, error) {
	metric, err := h.service.Get(req.Id)
	if err != nil {
		h.logger.Error("failed to get metric", zap.Error(err))
		return nil, err
	}

	if metric == nil {
		return nil, errors.New("metric not found")
	}

	resp := &proto.GetMetricResponse{
		Id:   metric.Name,
		Type: metric.Type,
	}

	if metric.Type == models.GaugeType && metric.Value != nil {
		resp.Value = *metric.Value
	} else if metric.Type == models.CounterType && metric.Delta != nil {
		resp.Delta = *metric.Delta
	}

	return resp, nil
}
