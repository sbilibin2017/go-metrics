package handlers_test

import (
	"context"
	"errors"
	"go-metrics/internal/handlers"
	"go-metrics/internal/responses"
	"go-metrics/internal/validation"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	tests := []struct {
		name               string
		mockResponse       *responses.MetricUpdatePathResponse
		mockError          error
		expectedStatusCode int
	}{
		{
			name:               "Success",
			mockResponse:       &responses.MetricUpdatePathResponse{},
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Error EmptyName",
			mockResponse:       nil,
			mockError:          validation.ErrEmptyName,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Error InvalidType",
			mockResponse:       nil,
			mockError:          validation.ErrInvalidType,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Error EmptyValue",
			mockResponse:       nil,
			mockError:          validation.ErrEmptyValue,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Error InternalServerError",
			mockResponse:       nil,
			mockError:          errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := handlers.NewMockMetricUpdatePathUsecase(ctrl)
			mockUsecase.EXPECT().Execute(context.Background(), gomock.Any()).Return(tt.mockResponse, tt.mockError)

			req := httptest.NewRequest(http.MethodPost, "/update/{type}/{name}/{value}", nil)
			req = req.WithContext(context.Background())
			rr := httptest.NewRecorder()
			handler := handlers.MetricUpdatePathHandler(mockUsecase)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

		})
	}
}
