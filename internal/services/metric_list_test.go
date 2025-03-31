package services_test

import (
	"context"
	"errors"
	"testing"

	"go-metrics/internal/domain"
	"go-metrics/internal/services"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test setup function to initialize mock repo and controller
func setup(t *testing.T) (*gomock.Controller, *services.MockMetricListFindRepository, *services.MetricListService) {
	ctrl := gomock.NewController(t)
	mockRepo := services.NewMockMetricListFindRepository(ctrl)
	service := services.NewMetricListService(mockRepo)
	return ctrl, mockRepo, service
}

// Test 1: Successfully retrieving the metric list
func TestMetricListService_List_Success(t *testing.T) {
	ctrl, mockRepo, service := setup(t)
	defer ctrl.Finish()

	// Mock expected behavior
	mockRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		domain.MetricID{ID: "metric-1", MType: "counter"}: {ID: "metric-1", MType: "counter", Delta: nil, Value: nil},
		domain.MetricID{ID: "metric-2", MType: "gauge"}:   {ID: "metric-2", MType: "gauge", Delta: nil, Value: nil},
	}, nil).Times(1)

	// Call the List method
	result, err := service.List(context.Background())

	// Assert no error and check result
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "metric-1", result[0].ID)
	assert.Equal(t, "metric-2", result[1].ID)
}

// Test 2: Error from repository (internal error)
func TestMetricListService_List_FindRepositoryError(t *testing.T) {
	ctrl, mockRepo, service := setup(t)
	defer ctrl.Finish()

	// Mock Find to return an error
	mockRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("repository error")).Times(1)

	// Call the List method
	result, err := service.List(context.Background())

	// Assert that the error is wrapped with ErrMetricListInternal
	assert.Equal(t, services.ErrMetricListInternal, err)
	assert.Nil(t, result)
}

// Test 3: Empty result from repository (No metrics found)
func TestMetricListService_List_EmptyResult(t *testing.T) {
	ctrl, mockRepo, service := setup(t)
	defer ctrl.Finish()

	// Mock Find to return an empty map
	mockRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{}, nil).Times(1)

	// Call the List method
	result, err := service.List(context.Background())

	// Assert no error and that the result is an empty slice
	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

// Test 4: List method returns metrics sorted by ID
func TestMetricListService_List_SortedResult(t *testing.T) {
	ctrl, mockRepo, service := setup(t)
	defer ctrl.Finish()

	// Mock Find to return unordered metrics
	mockRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(map[domain.MetricID]*domain.Metric{
		domain.MetricID{ID: "metric-2", MType: "gauge"}:   {ID: "metric-2", MType: "gauge", Delta: nil, Value: nil},
		domain.MetricID{ID: "metric-1", MType: "counter"}: {ID: "metric-1", MType: "counter", Delta: nil, Value: nil},
	}, nil).Times(1)

	// Call the List method
	result, err := service.List(context.Background())

	// Assert no error and check that the result is sorted by ID
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "metric-1", result[0].ID)
	assert.Equal(t, "metric-2", result[1].ID)
}
