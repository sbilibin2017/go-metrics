package usecases

import (
	"context"
	"errors"
	"testing"

	"go-metrics/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockService)
	ctx := context.Background()
	metric := &domain.Metric{}
	reqMock := NewMockMetricUpdatePathRequest(ctrl)
	reqMock.EXPECT().Validate().Return(nil)
	reqMock.EXPECT().ToDomain().Return(metric, nil)
	mockService.EXPECT().Update(ctx, []*domain.Metric{metric}).Return([]*domain.Metric{metric}, nil)
	response, err := usecase.Execute(ctx, reqMock)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestMetricUpdatePathUsecase_Execute_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockService)
	ctx := context.Background()
	reqMock := NewMockMetricUpdatePathRequest(ctrl)
	reqMock.EXPECT().Validate().Return(errors.New("validation error"))
	response, err := usecase.Execute(ctx, reqMock)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestMetricUpdatePathUsecase_Execute_UpdateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockService)
	ctx := context.Background()
	reqMock := NewMockMetricUpdatePathRequest(ctrl)
	metric := &domain.Metric{}
	reqMock.EXPECT().Validate().Return(nil)
	reqMock.EXPECT().ToDomain().Return(metric, nil)
	mockService.EXPECT().Update(ctx, []*domain.Metric{metric}).Return(nil, errors.New("update error"))
	response, err := usecase.Execute(ctx, reqMock)
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestMetricUpdatePathUsecase_Execute_ToDomainError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := NewMetricUpdatePathUsecase(mockService)
	ctx := context.Background()
	reqMock := NewMockMetricUpdatePathRequest(ctrl)
	reqMock.EXPECT().Validate().Return(nil)
	reqMock.EXPECT().ToDomain().Return(nil, errors.New("ToDomain error"))
	response, err := usecase.Execute(ctx, reqMock)
	assert.Error(t, err)
	assert.Nil(t, response)
}
