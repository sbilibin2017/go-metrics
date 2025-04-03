package handlers

import (
	"encoding/json"
	e "errors"
	"go-metrics/internal/usecases"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)
	mockResponse := usecases.MetricUpdatePathResponse("test")
	mockUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(&mockResponse, nil).
		Times(1)
	r := chi.NewRouter()
	r.Get("/metrics/update/{type}/{name}/{value}", MetricUpdatePathHandler(mockUsecase))
	req := httptest.NewRequest(http.MethodGet, "/metrics/update/type/name/value", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	_, err := json.Marshal(mockResponse)
	require.NoError(t, err)

}

func TestMetricUpdatePathHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathUsecase
	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)

	// Set up the expectation for the Execute method call to return an error
	mockUsecase.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(nil, e.New("some error")).
		Times(1)

	// Create a new chi router and set up the handler
	r := chi.NewRouter()
	r.Get("/metrics/update/{type}/{name}/{value}", MetricUpdatePathHandler(mockUsecase))

	// Create a new HTTP request with URL parameters matching the route
	req := httptest.NewRequest(http.MethodGet, "/metrics/update/type/name/value", nil)

	// Create a new ResponseRecorder to capture the HTTP response
	rr := httptest.NewRecorder()

	// Call the router with the request
	r.ServeHTTP(rr, req)

	// Verify the response status code and body for error case
	require.Equal(t, http.StatusInternalServerError, rr.Code)

}
