package services_test

import (
	"context"
	e "errors"
	"go-metrics/internal/domain"
	"go-metrics/internal/errors"
	"go-metrics/internal/services"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFindRepo := services.NewMockMetricListFindRepository(ctrl)
	metricsMap := map[domain.MetricID]*domain.Metric{
		{ID: "2", Type: domain.Counter}: {MetricID: domain.MetricID{ID: "2", Type: domain.Counter}, Value: new(float64)},
		{ID: "1", Type: domain.Counter}: {MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Value: new(float64)},
	}
	expectedMetrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Value: new(float64)},
		{MetricID: domain.MetricID{ID: "2", Type: domain.Counter}, Value: new(float64)},
	}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(metricsMap, nil).Times(1)
	service := services.NewMetricListService(mockFindRepo)
	result, err := service.List(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
}

func TestList_FindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFindRepo := services.NewMockMetricListFindRepository(ctrl)
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, e.New("find error")).Times(1)
	service := services.NewMetricListService(mockFindRepo)
	result, err := service.List(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, errors.ErrMetricListInternal.Error())
}
