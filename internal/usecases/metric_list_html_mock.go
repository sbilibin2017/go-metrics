// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/usecases/metric_list_html.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	domain "go-metrics/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricListHTMLService is a mock of MetricListHTMLService interface.
type MockMetricListHTMLService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricListHTMLServiceMockRecorder
}

// MockMetricListHTMLServiceMockRecorder is the mock recorder for MockMetricListHTMLService.
type MockMetricListHTMLServiceMockRecorder struct {
	mock *MockMetricListHTMLService
}

// NewMockMetricListHTMLService creates a new mock instance.
func NewMockMetricListHTMLService(ctrl *gomock.Controller) *MockMetricListHTMLService {
	mock := &MockMetricListHTMLService{ctrl: ctrl}
	mock.recorder = &MockMetricListHTMLServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricListHTMLService) EXPECT() *MockMetricListHTMLServiceMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockMetricListHTMLService) List(ctx context.Context) ([]*domain.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]*domain.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockMetricListHTMLServiceMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMetricListHTMLService)(nil).List), ctx)
}
