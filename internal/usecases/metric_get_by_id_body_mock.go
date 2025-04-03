// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/usecases/metric_get_by_id_body.go

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	domain "go-metrics/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricGetByIDBodyService is a mock of MetricGetByIDBodyService interface.
type MockMetricGetByIDBodyService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricGetByIDBodyServiceMockRecorder
}

// MockMetricGetByIDBodyServiceMockRecorder is the mock recorder for MockMetricGetByIDBodyService.
type MockMetricGetByIDBodyServiceMockRecorder struct {
	mock *MockMetricGetByIDBodyService
}

// NewMockMetricGetByIDBodyService creates a new mock instance.
func NewMockMetricGetByIDBodyService(ctrl *gomock.Controller) *MockMetricGetByIDBodyService {
	mock := &MockMetricGetByIDBodyService{ctrl: ctrl}
	mock.recorder = &MockMetricGetByIDBodyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricGetByIDBodyService) EXPECT() *MockMetricGetByIDBodyServiceMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockMetricGetByIDBodyService) GetByID(ctx context.Context, id *domain.MetricID) (*domain.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockMetricGetByIDBodyServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMetricGetByIDBodyService)(nil).GetByID), ctx, id)
}
