package services

import (
	"context"
	"go-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdateService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSaveRepo := NewMockMetricUpdateSaveBatchRepository(ctrl)
	mockFindRepo := NewMockMetricUpdateFindBatchRepository(ctrl)
	mockUnitOfWork := NewMockUnitOfWork(ctrl)
	service := NewMetricUpdateService(mockSaveRepo, mockFindRepo, mockUnitOfWork)

	metric := &domain.Metric{
		ID:    "metric1",
		MType: domain.Counter,
		Delta: new(int64),
	}
	*metric.Delta = 5

	existingMetrics := map[domain.MetricID]*domain.Metric{
		{ID: "metric1", MType: domain.Counter}: {
			ID:    "metric1",
			MType: domain.Counter,
			Delta: new(int64),
		},
	}
	*existingMetrics[domain.MetricID{ID: "metric1", MType: domain.Counter}].Delta = 10

	tests := []struct {
		name          string
		mockSetup     func()
		metrics       []*domain.Metric
		expectedError error
		expectedDelta int64
	}{
		{
			name: "should update metrics successfully",
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
					return operation()
				}).Times(1)
				mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(existingMetrics, nil).Times(1)
				mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			metrics:       []*domain.Metric{metric},
			expectedError: nil,
			expectedDelta: 15,
		},
		{
			name: "should return error if Find fails",
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
					return operation()
				}).Times(1)
				mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, ErrMetricIsNotUpdated).Times(1)
			},
			metrics:       []*domain.Metric{metric},
			expectedError: ErrMetricIsNotUpdated,
			expectedDelta: 0,
		},
		{
			name: "should return error if Save fails",
			mockSetup: func() {
				mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func() error) error {
					return operation()
				}).Times(1)
				mockFindRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(existingMetrics, nil).Times(1)
				mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(ErrMetricIsNotUpdated).Times(1)
			},
			metrics:       []*domain.Metric{metric},
			expectedError: ErrMetricIsNotUpdated,
			expectedDelta: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			updatedMetrics, err := service.Update(context.Background(), tt.metrics)
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, len(updatedMetrics))
				assert.Equal(t, tt.expectedDelta, *updatedMetrics[0].Delta)
			}
		})
	}
}
