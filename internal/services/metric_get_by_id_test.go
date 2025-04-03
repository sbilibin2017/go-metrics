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

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFindRepo := services.NewMockMetricGetByIDFindRepository(ctrl)
	metricID := &domain.MetricID{ID: "1", Type: domain.Counter}
	expectedMetric := &domain.Metric{MetricID: *metricID, Value: new(float64)}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		*metricID: expectedMetric,
	}, nil).Times(1)
	service := services.NewMetricGetByIDService(mockFindRepo)
	result, err := service.GetByID(context.Background(), metricID)
	require.NoError(t, err)
	assert.Equal(t, expectedMetric, result)
}

func TestGetByID_MetricNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFindRepo := services.NewMockMetricGetByIDFindRepository(ctrl)
	metricID := &domain.MetricID{ID: "1", Type: domain.Counter}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{}, nil).Times(1)
	service := services.NewMetricGetByIDService(mockFindRepo)
	result, err := service.GetByID(context.Background(), metricID)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, errors.ErrMetricNotFound.Error())
}

func TestGetByID_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFindRepo := services.NewMockMetricGetByIDFindRepository(ctrl)
	metricID := &domain.MetricID{ID: "1", Type: domain.Counter}
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, e.New("db error")).Times(1)
	service := services.NewMetricGetByIDService(mockFindRepo)
	result, err := service.GetByID(context.Background(), metricID)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, errors.ErrMetricGetByIDInternal.Error())
}
