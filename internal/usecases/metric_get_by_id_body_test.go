package usecases

import (
	"context"
	"errors"
	"testing"

	"go-metrics/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetByIDBodyUsecase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		req          *MetricGetByIDBodyRequest
		mock         func(mockService *MockMetricGetByIDBodyService, req *MetricGetByIDBodyRequest)
		expectErr    bool
		expectedResp *MetricGetByIDBodyResponse
	}{
		{
			name: "success case - counter",
			req: &MetricGetByIDBodyRequest{
				Type: string(domain.Counter),
				ID:   "test_counter",
			},
			mock: func(mockService *MockMetricGetByIDBodyService, req *MetricGetByIDBodyRequest) {
				metricID := ConvertMetricGetByIDBodyRequestToDomain(req)
				delta := int64(100)
				mockService.EXPECT().
					GetByID(gomock.Any(), metricID).
					Return(&domain.Metric{
						MetricID: *metricID,
						Delta:    &delta,
					}, nil)
			},
			expectErr: false,
			expectedResp: &MetricGetByIDBodyResponse{
				ID:    "test_counter",
				Type:  string(domain.Counter),
				Delta: int64Ptr(100),
				Value: nil,
			},
		},
		{
			name: "success case - gauge",
			req: &MetricGetByIDBodyRequest{
				Type: string(domain.Gauge),
				ID:   "test_gauge",
			},
			mock: func(mockService *MockMetricGetByIDBodyService, req *MetricGetByIDBodyRequest) {
				metricID := ConvertMetricGetByIDBodyRequestToDomain(req)
				value := 123.456
				mockService.EXPECT().
					GetByID(gomock.Any(), metricID).
					Return(&domain.Metric{
						MetricID: *metricID,
						Value:    &value,
					}, nil)
			},
			expectErr: false,
			expectedResp: &MetricGetByIDBodyResponse{
				ID:    "test_gauge",
				Type:  string(domain.Gauge),
				Delta: nil,
				Value: float64Ptr(123.456),
			},
		},
		{
			name: "validation error - empty ID",
			req: &MetricGetByIDBodyRequest{
				Type: string(domain.Counter),
				ID:   "",
			},
			mock:         nil,
			expectErr:    true,
			expectedResp: nil,
		},
		{
			name: "validation error - invalid type",
			req: &MetricGetByIDBodyRequest{
				Type: "invalid_type",
				ID:   "test_metric",
			},
			mock:         nil,
			expectErr:    true,
			expectedResp: nil,
		},
		{
			name: "service error - metric not found",
			req: &MetricGetByIDBodyRequest{
				Type: string(domain.Counter),
				ID:   "nonexistent_metric",
			},
			mock: func(mockService *MockMetricGetByIDBodyService, req *MetricGetByIDBodyRequest) {
				metricID := ConvertMetricGetByIDBodyRequestToDomain(req)
				mockService.EXPECT().
					GetByID(gomock.Any(), metricID).
					Return(nil, errors.New("metric not found"))
			},
			expectErr:    true,
			expectedResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockMetricGetByIDBodyService(ctrl)
			usecase := &MetricGetByIDBodyUsecase{svc: mockService}

			if tt.mock != nil {
				tt.mock(mockService, tt.req)
			}

			resp, err := usecase.Execute(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
