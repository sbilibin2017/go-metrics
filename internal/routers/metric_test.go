package routers

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go-metrics/internal/middlewares"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Фиктивная реализация Config
type mockConfig struct {
	key string
}

func (m *mockConfig) GetKey() string {
	return m.key
}

// Генерация HMAC-хеша
func computeHMAC(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Фиктивный обработчик
func mockHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "success"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func TestMiddlewareStack(t *testing.T) {
	secretKey := "supersecret"
	cfg := &mockConfig{key: secretKey}
	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)
	r.Use(middlewares.HMACMiddleware(cfg))
	r.Post("/test", mockHandler)
	requestData := map[string]string{"metric": "cpu", "value": "90"}
	jsonData, _ := json.Marshal(requestData)
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(jsonData)
	assert.NoError(t, err)
	gzipWriter.Close()
	req := httptest.NewRequest("POST", "/test", &compressedData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("HashSHA256", computeHMAC(jsonData, secretKey))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Должен быть статус 200 OK")
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json", "Content-Type должен быть application/json")
	expectedResponse := `{"message":"success"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String(), "Тело ответа должно быть корректным JSON")
}
