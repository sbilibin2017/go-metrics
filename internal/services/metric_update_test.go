package services_test

import (
	"context"
	"database/sql"
	e "errors"
	"go-metrics/internal/domain"
	"go-metrics/internal/errors"
	"go-metrics/internal/services"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_SuccessfulUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := services.NewMockMetricUpdateSaveRepository(ctrl)
	mockFindRepo := services.NewMockMetricUpdateFindRepository(ctrl)
	mockUnitOfWork := services.NewMockUnitOfWork(ctrl)
	metrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}
	expectedMetrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		{ID: "1", Type: domain.Counter}: {MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}, nil).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), metrics).Return(nil).Times(1)
	service := services.NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockUnitOfWork)
	result, err := service.Update(context.Background(), metrics)
	require.NoError(t, err)
	assert.Equal(t, expectedMetrics, result)
}

func TestUpdate_ExistingMetricsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := services.NewMockMetricUpdateSaveRepository(ctrl)
	mockFindRepo := services.NewMockMetricUpdateFindRepository(ctrl)
	mockUnitOfWork := services.NewMockUnitOfWork(ctrl)
	metrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), metrics).Return(nil).Times(1)
	service := services.NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockUnitOfWork)
	result, err := service.Update(context.Background(), metrics)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
}

func TestUpdate_Failure_FindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := services.NewMockMetricUpdateSaveRepository(ctrl)
	mockFindRepo := services.NewMockMetricUpdateFindRepository(ctrl)
	mockUnitOfWork := services.NewMockUnitOfWork(ctrl)
	metrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, e.New("find error")).Times(1)
	service := services.NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockUnitOfWork)
	result, err := service.Update(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "find error")
}

func TestUpdate_Failure_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSaveRepo := services.NewMockMetricUpdateSaveRepository(ctrl)
	mockFindRepo := services.NewMockMetricUpdateFindRepository(ctrl)
	mockUnitOfWork := services.NewMockUnitOfWork(ctrl)
	metrics := []*domain.Metric{
		{MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}
	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)
	mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		{ID: "1", Type: domain.Counter}: {MetricID: domain.MetricID{ID: "1", Type: domain.Counter}, Delta: new(int64)},
	}, nil).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), metrics).Return(e.New("save error")).Times(1)
	service := services.NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockUnitOfWork)
	result, err := service.Update(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, errors.ErrMetricIsNotUpdated.Error())
}
