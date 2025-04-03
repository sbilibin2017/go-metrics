package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestNewMetricRouter(t *testing.T) {
	h1 := new(mockHandler)
	h2 := new(mockHandler)
	h3 := new(mockHandler)
	h4 := new(mockHandler)
	h5 := new(mockHandler)
	h6 := new(mockHandler)

	h1.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	h2.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	h3.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	h4.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	h5.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	h6.On("ServeHTTP", mock.Anything, mock.Anything).Return()

	r := NewMetricRouter(h1.ServeHTTP, h2.ServeHTTP, h3.ServeHTTP, h4.ServeHTTP, h5.ServeHTTP, h6.ServeHTTP)

	tests := []struct {
		method  string
		url     string
		handler *mockHandler
	}{
		{"POST", "/update/metric/metricName/100", h1},
		{"POST", "/update/", h2},
		{"POST", "/updates/", h3},
		{"GET", "/value/metric/metricName", h4},
		{"GET", "/", h6},        // Исправлено (было h5)
		{"POST", "/value/", h5}, // Исправлено (было h6)
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			tt.handler.AssertNumberOfCalls(t, "ServeHTTP", 1)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}

func TestMiddlewares(t *testing.T) {
	h := new(mockHandler)
	h.On("ServeHTTP", mock.Anything, mock.Anything).Return()
	r := NewMetricRouter(h.ServeHTTP, h.ServeHTTP, h.ServeHTTP, h.ServeHTTP, h.ServeHTTP, h.ServeHTTP)
	req, err := http.NewRequest("POST", "/update/metric/metricName/100", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	h.AssertNumberOfCalls(t, "ServeHTTP", 1)
	assert.Equal(t, http.StatusOK, rr.Code)
}
