package usecases

import (
	"context"
	"errors"
	"testing"

	"go-metrics/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricGetByIDPathUsecase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		req          *MetricGetByIDPathRequest
		mock         func(mockService *MockMetricGetByIDPathService, req *MetricGetByIDPathRequest)
		expectErr    bool
		expectedResp *MetricGetByIDPathResponse
	}{
		{
			name: "success case - counter",
			req: &MetricGetByIDPathRequest{
				Type: string(domain.Counter),
				Name: "test_counter",
			},
			mock: func(mockService *MockMetricGetByIDPathService, req *MetricGetByIDPathRequest) {
				metricID := ConvertMetricGetByIDPathRequestToDomain(req)
				delta := int64(100)
				mockService.EXPECT().
					GetByID(gomock.Any(), metricID).
					Return(&domain.Metric{
						MetricID: *metricID,
						Delta:    &delta,
					}, nil)
			},
			expectErr:    false,
			expectedResp: toMetricGetByIDPathResponse("100"),
		},
		{
			name: "success case - gauge",
			req: &MetricGetByIDPathRequest{
				Type: string(domain.Gauge),
				Name: "test_gauge",
			},
			mock: func(mockService *MockMetricGetByIDPathService, req *MetricGetByIDPathRequest) {
				metricID := ConvertMetricGetByIDPathRequestToDomain(req)
				value := 123.456
				mockService.EXPECT().
					GetByID(gomock.Any(), metricID).
					Return(&domain.Metric{
						MetricID: *metricID,
						Value:    &value,
					}, nil)
			},
			expectErr:    false,
			expectedResp: toMetricGetByIDPathResponse("123.456"),
		},
		{
			name: "validation error - empty name",
			req: &MetricGetByIDPathRequest{
				Type: string(domain.Counter),
				Name: "",
			},
			mock:         nil,
			expectErr:    true,
			expectedResp: nil,
		},
		{
			name: "validation error - invalid type",
			req: &MetricGetByIDPathRequest{
				Type: "invalid_type",
				Name: "test_metric",
			},
			mock:         nil,
			expectErr:    true,
			expectedResp: nil,
		},
		{
			name: "service error - metric not found",
			req: &MetricGetByIDPathRequest{
				Type: string(domain.Counter),
				Name: "nonexistent_metric",
			},
			mock: func(mockService *MockMetricGetByIDPathService, req *MetricGetByIDPathRequest) {
				metricID := ConvertMetricGetByIDPathRequestToDomain(req)
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

			mockService := NewMockMetricGetByIDPathService(ctrl)
			usecase := &MetricGetByIDPathUsecase{svc: mockService}

			if tt.mock != nil {
				tt.mock(mockService, tt.req)
			}

			resp, err := usecase.Execute(context.Background(), tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedResp, *resp)
			}
		})
	}
}

// Вспомогательная функция для конвертации строки в указатель MetricGetByIDPathResponse
func toMetricGetByIDPathResponse(value string) *MetricGetByIDPathResponse {
	resp := MetricGetByIDPathResponse(value)
	return &resp
}
