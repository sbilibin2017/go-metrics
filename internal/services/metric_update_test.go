package services_test

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
	"go-metrics/internal/services"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := services.NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := services.NewMockMetricUpdateFindBatchRepository(ctrl)
	uow := services.NewMockUnitOfWork(ctrl)
	service := services.NewMetricUpdateService(saveRepo, findRepo, uow)
	metrics := []*domain.Metric{
		{ID: "1", Type: domain.Counter, Delta: new(int64)},
		{ID: "2", Type: domain.Counter, Delta: new(int64)},
	}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		{ID: "1", Type: domain.Counter}: {ID: "1", Type: domain.Counter, Delta: new(int64)},
	}, nil)
	saveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
	uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
		return operation()
	})
	updatedMetrics, err := service.Update(context.Background(), metrics)
	require.NoError(t, err)
	assert.Equal(t, metrics, updatedMetrics)
}

func TestMetricUpdateService_Update_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := services.NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := services.NewMockMetricUpdateFindBatchRepository(ctrl)
	uow := services.NewMockUnitOfWork(ctrl)
	service := services.NewMetricUpdateService(saveRepo, findRepo, uow)
	metrics := []*domain.Metric{
		{ID: "1", Type: domain.Counter, Delta: new(int64)},
		{ID: "2", Type: domain.Counter, Delta: new(int64)},
	}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		{ID: "1", Type: domain.Counter}: {ID: "1", Type: domain.Counter, Delta: new(int64)},
	}, nil)
	saveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save failed"))
	uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
		return operation()
	})
	updatedMetrics, err := service.Update(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, services.ErrMetricIsNotUpdated, err)
}

func TestMetricUpdateService_Update_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	saveRepo := services.NewMockMetricUpdateSaveBatchRepository(ctrl)
	findRepo := services.NewMockMetricUpdateFindBatchRepository(ctrl)
	uow := services.NewMockUnitOfWork(ctrl)
	service := services.NewMetricUpdateService(saveRepo, findRepo, uow)
	metrics := []*domain.Metric{
		{ID: "1", Type: domain.Counter, Delta: new(int64)},
		{ID: "2", Type: domain.Counter, Delta: new(int64)},
	}
	findRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("find failed"))
	uow.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
		return operation()
	})
	updatedMetrics, err := service.Update(context.Background(), metrics)
	require.Error(t, err)
	assert.Nil(t, updatedMetrics)
	assert.Equal(t, services.ErrMetricIsNotUpdated, err)
}
