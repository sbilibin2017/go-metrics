// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/handlers/metric_update_path.go

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	usecases "go-metrics/internal/usecases"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricUpdatePathUsecase is a mock of MetricUpdatePathUsecase interface.
type MockMetricUpdatePathUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockMetricUpdatePathUsecaseMockRecorder
}

// MockMetricUpdatePathUsecaseMockRecorder is the mock recorder for MockMetricUpdatePathUsecase.
type MockMetricUpdatePathUsecaseMockRecorder struct {
	mock *MockMetricUpdatePathUsecase
}

// NewMockMetricUpdatePathUsecase creates a new mock instance.
func NewMockMetricUpdatePathUsecase(ctrl *gomock.Controller) *MockMetricUpdatePathUsecase {
	mock := &MockMetricUpdatePathUsecase{ctrl: ctrl}
	mock.recorder = &MockMetricUpdatePathUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricUpdatePathUsecase) EXPECT() *MockMetricUpdatePathUsecaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockMetricUpdatePathUsecase) Execute(ctx context.Context, req *usecases.MetricUpdatePathRequest) (*usecases.MetricUpdatePathResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, req)
	ret0, _ := ret[0].(*usecases.MetricUpdatePathResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockMetricUpdatePathUsecaseMockRecorder) Execute(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockMetricUpdatePathUsecase)(nil).Execute), ctx, req)
}
