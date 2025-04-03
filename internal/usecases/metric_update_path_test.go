package usecases

import (
	"context"
	"testing"

	"go-metrics/internal/domain"
	"go-metrics/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathUsecase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		req          *MetricUpdatePathRequest
		mock         func(mockService *MockMetricUpdatePathService, req *MetricUpdatePathRequest)
		expectErr    error
		expectedResp *MetricUpdatePathResponse
	}{
		{
			name: "success case - counter",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "test_counter",
				Value: "123",
			},
			mock: func(mockService *MockMetricUpdatePathService, req *MetricUpdatePathRequest) {
				metric := ConvertMetricUpdatePathRequestToDomain(req)
				mockService.EXPECT().Update(gomock.Any(), []*domain.Metric{metric}).Return([]*domain.Metric{metric}, nil)
			},
			expectErr:    nil,
			expectedResp: NewMetricUpdatePathResponse(),
		},
		{
			name: "validation error - empty name",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "",
				Value: "123",
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidMetricID,
			expectedResp: nil,
		},
		{
			name: "invalid metric type",
			req: &MetricUpdatePathRequest{
				Type:  "invalid_type",
				Name:  "testgauge",
				Value: "123.456",
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidMetricType,
			expectedResp: nil,
		},
		{
			name: "invalid counter value",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "testcounter",
				Value: "abc",
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidCounterMetricValue,
			expectedResp: nil,
		},
		{
			name: "invalid gauge value",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "testgauge",
				Value: "abc",
			},
			mock:         nil,
			expectErr:    errors.ErrInvalidGaugeMetricValue,
			expectedResp: nil,
		},
		{
			name: "success case - gauge value conversion",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Gauge),
				Name:  "test_gauge",
				Value: "123.456",
			},
			mock: func(mockService *MockMetricUpdatePathService, req *MetricUpdatePathRequest) {
				metric := ConvertMetricUpdatePathRequestToDomain(req)
				mockService.EXPECT().Update(gomock.Any(), []*domain.Metric{metric}).Return([]*domain.Metric{metric}, nil)
			},
			expectErr:    nil,
			expectedResp: NewMetricUpdatePathResponse(),
		},
		{
			name: "service error",
			req: &MetricUpdatePathRequest{
				Type:  string(domain.Counter),
				Name:  "testcounter",
				Value: "123",
			},
			mock: func(mockService *MockMetricUpdatePathService, req *MetricUpdatePathRequest) {
				metric := ConvertMetricUpdatePathRequestToDomain(req)
				mockService.EXPECT().Update(gomock.Any(), []*domain.Metric{metric}).Return(nil, errors.ErrMetricIsNotUpdated)
			},
			expectErr:    errors.ErrMetricIsNotUpdated,
			expectedResp: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := NewMockMetricUpdatePathService(ctrl)
			usecase := &MetricUpdatePathUsecase{svc: mockService}
			if tt.mock != nil {
				tt.mock(mockService, tt.req)
			}
			resp, err := usecase.Execute(context.Background(), tt.req)
			if tt.expectErr != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Equal(t, tt.expectErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, *tt.expectedResp, *resp)
			}
		})
	}
}
