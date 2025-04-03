package usecases

import (
	"context"
	"errors"
	"go-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricListHTMLUsecase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricListHTMLService(ctrl)
	usecase := &MetricListHTMLUsecase{svc: mockService}

	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr error
		checkHTML   func(html string)
	}{
		{
			name: "success - generate HTML with metrics",
			mockSetup: func() {
				mockService.EXPECT().List(gomock.Any()).Return([]*domain.Metric{
					{MetricID: domain.MetricID{ID: "metric1", Type: domain.Counter}, Delta: ptrInt64(100)},
					{MetricID: domain.MetricID{ID: "metric2", Type: domain.Gauge}, Value: ptrFloat64(123.456)},
				}, nil)
			},
			expectedErr: nil,
			checkHTML: func(html string) {
				assert.Contains(t, html, "<table")
				assert.Contains(t, html, "<tr><th>ID</th><th>Value</th></tr>")
				assert.Contains(t, html, "metric1")
				assert.Contains(t, html, "100")
				assert.Contains(t, html, "metric2")
				assert.Contains(t, html, "123.456")
			},
		},
		{
			name: "error - service failure",
			mockSetup: func() {
				mockService.EXPECT().List(gomock.Any()).Return(nil, errors.New("service error"))
			},
			expectedErr: errors.New("service error"),
			checkHTML:   nil,
		},
		{
			name: "metric with missing values should show N/A",
			mockSetup: func() {
				mockService.EXPECT().List(gomock.Any()).Return([]*domain.Metric{
					{MetricID: domain.MetricID{ID: "metric3", Type: domain.Counter}, Delta: nil},
					{MetricID: domain.MetricID{ID: "metric4", Type: domain.Gauge}, Value: nil},
				}, nil)
			},
			expectedErr: nil,
			checkHTML: func(html string) {
				assert.Contains(t, html, "metric3")
				assert.Contains(t, html, "N/A")
				assert.Contains(t, html, "metric4")
				assert.Contains(t, html, "N/A")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := usecase.Execute(context.Background())

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if tt.checkHTML != nil {
					tt.checkHTML(resp.HTML)
				}
			}
		})
	}
}

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
