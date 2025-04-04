// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/services/metric_list.go

// Package services is a generated GoMock package.
package services

import (
	context "context"
	domain "go-metrics/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricListFindRepository is a mock of MetricListFindRepository interface.
type MockMetricListFindRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMetricListFindRepositoryMockRecorder
}

// MockMetricListFindRepositoryMockRecorder is the mock recorder for MockMetricListFindRepository.
type MockMetricListFindRepositoryMockRecorder struct {
	mock *MockMetricListFindRepository
}

// NewMockMetricListFindRepository creates a new mock instance.
func NewMockMetricListFindRepository(ctrl *gomock.Controller) *MockMetricListFindRepository {
	mock := &MockMetricListFindRepository{ctrl: ctrl}
	mock.recorder = &MockMetricListFindRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricListFindRepository) EXPECT() *MockMetricListFindRepositoryMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockMetricListFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, filters)
	ret0, _ := ret[0].(map[domain.MetricID]*domain.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockMetricListFindRepositoryMockRecorder) Find(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockMetricListFindRepository)(nil).Find), ctx, filters)
}
