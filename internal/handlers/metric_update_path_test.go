package handlers

import (
	"context"
	"errors"
	"go-metrics/internal/requests"
	"go-metrics/internal/responses"
	"go-metrics/internal/validation"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathHandler_ServeHTTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)
	handler := NewMetricUpdatePathHandler(mockUsecase)
	metricUpdateRequest := &requests.MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "requests",
		Value: "100",
	}
	mockResponse := &responses.MetricUpdatePathResponse{}
	mockUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Eq(metricUpdateRequest)).
		Return(mockResponse, nil)
	req := httptest.NewRequest(http.MethodPost, "/update/counter/value/100", nil)
	ps := httprouter.Params{
		{Key: "type", Value: "counter"},
		{Key: "name", Value: "requests"},
		{Key: "value", Value: "100"},
	}
	ctx := context.WithValue(req.Context(), httprouter.ParamsKey, ps)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestHandleMetricUpdatePathError_ErrEmptyName(t *testing.T) {
	rec := httptest.NewRecorder()
	err := validation.ErrEmptyName
	handleMetricUpdatePathError(rec, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleMetricUpdatePathError_ErrInvalidType(t *testing.T) {
	rec := httptest.NewRecorder()
	err := validation.ErrInvalidType
	handleMetricUpdatePathError(rec, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleMetricUpdatePathError_ErrEmptyValue(t *testing.T) {
	rec := httptest.NewRecorder()
	err := validation.ErrEmptyValue
	handleMetricUpdatePathError(rec, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleMetricUpdatePathError_InternalServerError(t *testing.T) {
	rec := httptest.NewRecorder()
	err := errors.New("some unknown error")
	handleMetricUpdatePathError(rec, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestMetricUpdatePathHandler_ServeHTTP_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)
	handler := NewMetricUpdatePathHandler(mockUsecase)
	req := httptest.NewRequest(http.MethodPost, "/metrics/update/path", nil)
	rec := httptest.NewRecorder()
	mockError := validation.ErrEmptyName
	mockUsecase.EXPECT().
		Execute(gomock.Eq(context.Background()), gomock.Any()).
		Return(nil, mockError)
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
