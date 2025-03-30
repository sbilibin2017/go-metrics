package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestRegisterMetricUpdatePathRouter(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Success POST Request",
			method:             http.MethodPost,
			url:                "/update/someType/someName/someValue",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Method (GET Request)",
			method:             http.MethodGet,
			url:                "/update/someType/someName/someValue",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "Invalid URL Path",
			method:             http.MethodPost,
			url:                "/update/invalidPath",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Missing Parameters",
			method:             http.MethodPost,
			url:                "/update/someType/someName",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Extra Slash in URL Path",
			method:             http.MethodPost,
			url:                "/update/someType/someName//",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			RegisterMetricUpdatePathRouter(r, mockHandler)
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
