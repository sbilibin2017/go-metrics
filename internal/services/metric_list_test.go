package services

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricListService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricListFindBatchRepository(ctrl)
	expectedMetrics := []*domain.Metric{
		{ID: "123", Type: domain.Counter},
		{ID: "124", Type: domain.Gauge},
	}
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{}).
		Return(map[domain.MetricID]*domain.Metric{
			{ID: "123", Type: domain.Counter}: expectedMetrics[0],
			{ID: "124", Type: domain.Gauge}:   expectedMetrics[1],
		}, nil).
		Times(1)
	service := NewMetricListService(mockRepo)
	result, err := service.List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
}

func TestMetricListService_List_ErrorFindingMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricListFindBatchRepository(ctrl)
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{}).
		Return(nil, errors.New("db error")).
		Times(1)
	service := NewMetricListService(mockRepo)
	result, err := service.List(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrMetricListInternal, err)
}

func TestMetricListService_List_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockMetricListFindBatchRepository(ctrl)
	mockRepo.EXPECT().Find(context.Background(), []domain.MetricID{}).
		Return(map[domain.MetricID]*domain.Metric{}, nil).
		Times(1)
	service := NewMetricListService(mockRepo)
	result, err := service.List(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}
