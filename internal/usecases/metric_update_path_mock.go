// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/usecases/metric_update_path.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	domain "go-metrics/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricUpdatePathService is a mock of MetricUpdatePathService interface.
type MockMetricUpdatePathService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricUpdatePathServiceMockRecorder
}

// MockMetricUpdatePathServiceMockRecorder is the mock recorder for MockMetricUpdatePathService.
type MockMetricUpdatePathServiceMockRecorder struct {
	mock *MockMetricUpdatePathService
}

// NewMockMetricUpdatePathService creates a new mock instance.
func NewMockMetricUpdatePathService(ctrl *gomock.Controller) *MockMetricUpdatePathService {
	mock := &MockMetricUpdatePathService{ctrl: ctrl}
	mock.recorder = &MockMetricUpdatePathServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricUpdatePathService) EXPECT() *MockMetricUpdatePathServiceMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockMetricUpdatePathService) Update(ctx context.Context, metrics []*domain.Metric) ([]*domain.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, metrics)
	ret0, _ := ret[0].([]*domain.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockMetricUpdatePathServiceMockRecorder) Update(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMetricUpdatePathService)(nil).Update), ctx, metrics)
}
