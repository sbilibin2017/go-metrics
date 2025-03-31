package services

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetByIDService_GetByID_MetricFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindRepository(ctrl)
	service := NewMetricGetByIDService(mockRepo)
	metricID := &domain.MetricID{ID: "metric1", MType: "gauge"}
	metric := &domain.Metric{ID: "metric1", MType: "gauge", Value: func() *float64 { v := 42.5; return &v }()}
	mockRepo.EXPECT().Find(context.Background(), []*domain.MetricID{metricID}).Return(map[domain.MetricID]*domain.Metric{*metricID: metric}, nil)
	result, err := service.GetByID(context.Background(), metricID)
	assert.NoError(t, err)
	assert.Equal(t, metric, result)
}

func TestMetricGetByIDService_GetByID_MetricNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindRepository(ctrl)
	service := NewMetricGetByIDService(mockRepo)
	metricID := &domain.MetricID{ID: "metric2", MType: "counter"}
	mockRepo.EXPECT().Find(context.Background(), []*domain.MetricID{metricID}).Return(map[domain.MetricID]*domain.Metric{}, nil)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Error(t, err)
	assert.Equal(t, ErrMetricNotFound, err)
	assert.Nil(t, result)
}

func TestMetricGetByIDService_GetByID_FindRepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricGetByIDFindRepository(ctrl)
	service := NewMetricGetByIDService(mockRepo)
	metricID := &domain.MetricID{ID: "metric3", MType: "gauge"}
	mockRepo.EXPECT().Find(context.Background(), []*domain.MetricID{metricID}).Return(nil, assert.AnError)
	result, err := service.GetByID(context.Background(), metricID)
	assert.Error(t, err)
	assert.Equal(t, ErrMetricGetByIDInternal, err)
	assert.Nil(t, result)
}
