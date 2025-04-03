package usecases

import (
	"context"
	"testing"

	"go-metrics/internal/domain"
	"go-metrics/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatesBodyUsecase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		req          []*MetricUpdateBodyRequest
		mock         func(mockService *MockMetricUpdatesBodyService, req []*MetricUpdateBodyRequest)
		expectErr    error
		expectedResp []*MetricUpdateBodyResponse
	}{
		{
			name: "success case - counter",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "test_counter",
					Type:  string(domain.Counter),
					Delta: int64Ptr(123),
				},
			},
			mock: func(mockService *MockMetricUpdatesBodyService, req []*MetricUpdateBodyRequest) {
				metrics := make([]*domain.Metric, len(req))
				for i, r := range req {
					metrics[i] = ConvertMetricUpdateBodyRequestToDomain(r)
				}
				mockService.EXPECT().Update(gomock.Any(), metrics).Return(metrics, nil)
			},
			expectErr:    nil,
			expectedResp: []*MetricUpdateBodyResponse{{ID: "test_counter", Type: "counter", Delta: int64Ptr(123)}},
		},
		{
			name: "validation error - empty id",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "",
					Type:  string(domain.Counter),
					Delta: int64Ptr(123),
				},
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidMetricID,
			expectedResp: nil,
		},
		{
			name: "invalid metric type",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "testgauge",
					Type:  "invalid_type", // Неверный тип метрики
					Delta: int64Ptr(123),
				},
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidMetricType,
			expectedResp: nil,
		},
		{
			name: "validation error - invalid counter delta",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "test_counter",
					Type:  string(domain.Counter),
					Delta: nil, // Необходимо передать значение для Delta
				},
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidCounterMetricValue,
			expectedResp: nil,
		},
		{
			name: "validation error - invalid gauge value",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "test_gauge",
					Type:  string(domain.Gauge),
					Value: nil, // Необходимо передать значение для Value
				},
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidGaugeMetricValue,
			expectedResp: nil,
		},
		{
			name: "success case - multiple metrics",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "test_counter1",
					Type:  string(domain.Counter),
					Delta: int64Ptr(123),
				},
				{
					ID:    "test_gauge1",
					Type:  string(domain.Gauge),
					Value: float64Ptr(456.78),
				},
			},
			mock: func(mockService *MockMetricUpdatesBodyService, req []*MetricUpdateBodyRequest) {
				metrics := make([]*domain.Metric, len(req))
				for i, r := range req {
					metrics[i] = ConvertMetricUpdateBodyRequestToDomain(r)
				}
				mockService.EXPECT().Update(gomock.Any(), metrics).Return(metrics, nil)
			},
			expectErr: nil,
			expectedResp: []*MetricUpdateBodyResponse{
				{ID: "test_counter1", Type: "counter", Delta: int64Ptr(123)},
				{ID: "test_gauge1", Type: "gauge", Value: float64Ptr(456.78)},
			},
		},
		{
			name: "service error",
			req: []*MetricUpdateBodyRequest{
				{
					ID:    "test_counter",
					Type:  string(domain.Counter),
					Delta: int64Ptr(123),
				},
			},
			mock: func(mockService *MockMetricUpdatesBodyService, req []*MetricUpdateBodyRequest) {
				metrics := make([]*domain.Metric, len(req))
				for i, r := range req {
					metrics[i] = ConvertMetricUpdateBodyRequestToDomain(r)
				}
				mockService.EXPECT().Update(gomock.Any(), metrics).Return(nil, errors.ErrMetricIsNotUpdated)
			},
			expectErr:    errors.ErrMetricIsNotUpdated,
			expectedResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockMetricUpdatesBodyService(ctrl)
			usecase := &MetricUpdatesBodyUsecase{svc: mockService}

			// Настройка мока
			if tt.mock != nil {
				tt.mock(mockService, tt.req)
			}

			// Вызов Execute
			resp, err := usecase.Execute(context.Background(), tt.req)

			// Проверки
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
