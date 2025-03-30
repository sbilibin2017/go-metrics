package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAddresser struct {
	mock.Mock
}

func (m *MockAddresser) GetAddress() string {
	args := m.Called()
	return args.String(0)
}

func TestNewServer(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)
	assert.NotNil(t, server)
}

func TestAddRouter(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)
	rtr := chi.NewRouter()
	rtr.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	})
	server.AddRouter(rtr)
	ts := httptest.NewServer(server.router)
	defer ts.Close()
	resp, err := ts.Client().Get(ts.URL + "/test")
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}
	defer resp.Body.Close() // Close the response body

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error while reading response body: %v", err)
	}
	assert.Equal(t, "Test", string(body))
}

func TestServerStartAndShutdown(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go func() {
		err := server.Start(ctx)
		assert.NoError(t, err)
	}()
	time.Sleep(500 * time.Millisecond)
	cancel()
	err := server.Start(ctx)
	assert.NoError(t, err)
}

func TestServerShutdownTimeout(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)
	if server.server == nil {
		t.Fatalf("Expected server to be initialized, but got nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shutdownCancel()
	go func() {
		err := server.Start(ctx)
		if err != nil {
			t.Errorf("Server start failed: %v", err)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	err := server.server.Shutdown(shutdownCtx)
	assert.Nil(t, err)
	if shutdownCtx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected context to be cancelled due to timeout, but got: %v", shutdownCtx.Err())
	}
}

type MockHTTPServer struct {
	ListenAndServeFunc func() error
	ShutdownFunc       func(ctx context.Context) error
}

func (m *MockHTTPServer) ListenAndServe() error {
	if m.ListenAndServeFunc != nil {
		return m.ListenAndServeFunc()
	}
	return nil // Default to no error
}

func (m *MockHTTPServer) Shutdown(ctx context.Context) error {
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc(ctx)
	}
	return nil // Default to no error
}

func TestServerStartError(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)

	mockServer := &MockHTTPServer{
		ShutdownFunc: func(ctx context.Context) error {
			return nil
		},
	}
	server.server = mockServer
	mockServer.ListenAndServeFunc = func() error {
		return errors.New("Server start error")
	}
	go func() {
		err := server.Start(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "Server start error", err.Error())
	}()
}

func TestServerShutdownError(t *testing.T) {
	mockAddr := "localhost:8080"
	mockAddresser := new(MockAddresser)
	mockAddresser.On("GetAddress").Return(mockAddr)
	server := NewServer(mockAddresser)
	mockServer := &MockHTTPServer{
		ListenAndServeFunc: func() error {
			return nil
		},
		ShutdownFunc: func(ctx context.Context) error {
			return errors.New("shutdown failed")
		},
	}
	server.server = mockServer
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := server.server.Shutdown(shutdownCtx)
	assert.Error(t, err)
}
