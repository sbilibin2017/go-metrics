package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestAddHandler(t *testing.T) {
	r := NewRouter()
	testHandler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Hello, World!"))
	}
	r.AddHandler(http.MethodGet, "/test", testHandler)
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServerHTTP().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, World!", rr.Body.String())
}

func TestAddMiddleware(t *testing.T) {
	r := NewRouter()
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "Middleware Applied")
			next.ServeHTTP(w, r)
		})
	}
	r.AddMiddleware(middleware)
	testHandler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Hello, Middleware!"))
	}
	r.AddHandler(http.MethodGet, "/middleware", testHandler)
	req, err := http.NewRequest(http.MethodGet, "/middleware", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServerHTTP().ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, Middleware!", rr.Body.String())
	assert.Equal(t, "Middleware Applied", rr.Header().Get("X-Custom-Header"))
}

func TestServerHTTPWithMultipleRoutesAndMiddleware(t *testing.T) {
	r := NewRouter()
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "Applied")
			next.ServeHTTP(w, r)
		})
	}
	r.AddMiddleware(middleware)
	handler1 := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Route 1"))
	}
	handler2 := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Route 2"))
	}
	r.AddHandler(http.MethodGet, "/route1", handler1)
	r.AddHandler(http.MethodGet, "/route2", handler2)
	req1, err := http.NewRequest(http.MethodGet, "/route1", nil)
	assert.NoError(t, err)
	rr1 := httptest.NewRecorder()
	r.ServerHTTP().ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)
	assert.Equal(t, "Route 1", rr1.Body.String())
	assert.Equal(t, "Applied", rr1.Header().Get("X-Middleware"))
	req2, err := http.NewRequest(http.MethodGet, "/route2", nil)
	assert.NoError(t, err)
	rr2 := httptest.NewRecorder()
	r.ServerHTTP().ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)
	assert.Equal(t, "Route 2", rr2.Body.String())
	assert.Equal(t, "Applied", rr2.Header().Get("X-Middleware"))
}

func TestGetRoutes(t *testing.T) {
	r := NewRouter()
	r.AddHandler(http.MethodGet, "/route1", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {})
	r.AddHandler(http.MethodPost, "/route2", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {})
	routes := r.GetRoutes()
	assert.Len(t, routes, 2)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/route1", routes[0].Path)
	assert.Equal(t, "POST", routes[1].Method)
	assert.Equal(t, "/route2", routes[1].Path)
}
