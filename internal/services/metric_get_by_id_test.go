package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetByIDService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindBatchRepository(ctrl)
	metricID := &domain.MetricID{ID: "123", Type: domain.Counter}
	expectedMetric := &domain.Metric{ID: "123", Type: domain.Counter, Delta: nil, Value: nil}
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{*metricID}).
		Return(map[domain.MetricID]*domain.Metric{*metricID: expectedMetric}, nil).
		Times(1)
	service := NewMetricGetByIDService(mockRepo)
	result, err := service.GetByID(context.Background(), metricID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMetric, result)
}

func TestMetricGetByIDService_GetByID_ErrorFindingMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindBatchRepository(ctrl)
	metricID := &domain.MetricID{ID: "123", Type: domain.Counter}
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{*metricID}).
		Return(nil, errors.New("db error")).
		Times(1)
	service := NewMetricGetByIDService(mockRepo)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrMetricGetByIDInternal, err)
}

func TestMetricGetByIDService_GetByID_MetricNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindBatchRepository(ctrl)
	metricID := &domain.MetricID{ID: "123", Type: domain.Counter}
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{*metricID}).
		Return(map[domain.MetricID]*domain.Metric{}, nil).
		Times(1)
	service := NewMetricGetByIDService(mockRepo)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrMetricNotFound, err)
}
