package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/Zrossiz/go-metrics/internal/proto"
	"github.com/Zrossiz/go-metrics/internal/server/dto"
	"github.com/Zrossiz/go-metrics/internal/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockMetricService struct {
	mock.Mock
}

func (m *MockMetricService) Create(body dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockMetricService) Get(name string) (*models.Metric, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Metric), args.Error(1)
}

func (m *MockMetricService) GetAll() ([]models.Metric, error) {
	args := m.Called()
	return args.Get(0).([]models.Metric), args.Error(1)
}

func (m *MockMetricService) GetStringValueMetric(name string) (string, error) {
	args := m.Called(name)
	return args.String(0), args.Error(1)
}

func (m *MockMetricService) PingDB() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMetricService) SetBatch(body []dto.PostMetricDto) error {
	args := m.Called(body)
	return args.Error(0)
}

func TestUpdateParam_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("Create", mock.Anything).Return(nil)

	req := &proto.PostMetricRequest{
		Id:    "test_id",
		Type:  "gauge",
		Value: 10.5,
	}

	resp, err := handler.UpdateParam(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestUpdateParam_ErrorOnCreate(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("Create", mock.Anything).Return(errors.New("database error"))

	req := &proto.PostMetricRequest{
		Id:    "test_id",
		Type:  "gauge",
		Value: 10.5,
	}

	resp, err := handler.UpdateParam(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockService.AssertExpectations(t)
}

func TestUpdateJSON_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("Create", mock.Anything).Return(nil)

	req := &proto.PostMetricRequest{
		Id:    "test_id",
		Type:  "gauge",
		Value: 10.5,
	}

	resp, err := handler.UpdateJSON(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestUpdateBatchJSON_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("SetBatch", mock.Anything).Return(nil)

	req := &proto.BatchPostMetricRequest{
		Metrics: []*proto.PostMetricRequest{
			{
				Id:    "metric1",
				Type:  "gauge",
				Value: 100.5,
			},
			{
				Id:    "metric2",
				Type:  "counter",
				Delta: 5,
			},
		},
	}

	resp, err := handler.UpdateBatchJSON(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

func TestUpdateBatchJSON_ErrorOnSetBatch(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("SetBatch", mock.Anything).Return(errors.New("batch error"))

	req := &proto.BatchPostMetricRequest{
		Metrics: []*proto.PostMetricRequest{
			{
				Id:    "metric1",
				Type:  "gauge",
				Value: 100.5,
			},
			{
				Id:    "metric2",
				Type:  "counter",
				Delta: 5,
			},
		},
	}

	resp, err := handler.UpdateBatchJSON(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockService.AssertExpectations(t)
}

func TestGetMetric_Success(t *testing.T) {
	mockService := new(MockMetricService)
	logger, _ := zap.NewDevelopment()
	handler := NewGrpcHandler(mockService, logger)

	mockService.On("Get", "test_id").Return(&models.Metric{
		Name:  "test_id",
		Type:  "gauge",
		Value: &[]float64{10.5}[0],
	}, nil)

	req := &proto.GetMetricRequest{
		Id: "test_id",
	}

	resp, err := handler.GetMetric(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test_id", resp.Id)
	assert.Equal(t, "gauge", resp.Type)
	assert.Equal(t, 10.5, resp.Value)

	mockService.AssertExpectations(t)
}
