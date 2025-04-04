// Code generated by MockGen. DO NOT EDIT.
// Source: /home/sergey/Go/tmp2/go-metrics/internal/services/metric_get_by_id.go

// Package services is a generated GoMock package.
package services

import (
	context "context"
	domain "go-metrics/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricGetByIDFindRepository is a mock of MetricGetByIDFindRepository interface.
type MockMetricGetByIDFindRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMetricGetByIDFindRepositoryMockRecorder
}

// MockMetricGetByIDFindRepositoryMockRecorder is the mock recorder for MockMetricGetByIDFindRepository.
type MockMetricGetByIDFindRepositoryMockRecorder struct {
	mock *MockMetricGetByIDFindRepository
}

// NewMockMetricGetByIDFindRepository creates a new mock instance.
func NewMockMetricGetByIDFindRepository(ctrl *gomock.Controller) *MockMetricGetByIDFindRepository {
	mock := &MockMetricGetByIDFindRepository{ctrl: ctrl}
	mock.recorder = &MockMetricGetByIDFindRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricGetByIDFindRepository) EXPECT() *MockMetricGetByIDFindRepositoryMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockMetricGetByIDFindRepository) Find(ctx context.Context, filters []*domain.MetricID) (map[domain.MetricID]*domain.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, filters)
	ret0, _ := ret[0].(map[domain.MetricID]*domain.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockMetricGetByIDFindRepositoryMockRecorder) Find(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockMetricGetByIDFindRepository)(nil).Find), ctx, filters)
}
