package middlewares_test

import (
	"bytes"
	"compress/gzip"
	"go-metrics/internal/middlewares"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressRequestBody(t *testing.T) {
	t.Run("should decompress request body if Content-Encoding is gzip", func(t *testing.T) {
		inputData := "this is a test body"
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		_, err := gzipWriter.Write([]byte(inputData))
		require.NoError(t, err)
		err = gzipWriter.Close()
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/", &buf)
		req.Header.Set("Content-Encoding", "gzip")

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, inputData, string(body))
		})

		rr := httptest.NewRecorder()
		handler := middlewares.GzipMiddleware(mockHandler)
		handler.ServeHTTP(rr, req)
	})
}

func TestCompressResponse(t *testing.T) {
	t.Run("should compress response if Accept-Encoding contains gzip", func(t *testing.T) {
		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("response body"))
			require.NoError(t, err)
		})

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		rr := httptest.NewRecorder()
		handler := middlewares.GzipMiddleware(mockHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

		gzipReader, err := gzip.NewReader(rr.Body)
		require.NoError(t, err)
		defer gzipReader.Close()
		decompressedData, err := io.ReadAll(gzipReader)
		require.NoError(t, err)

		assert.Equal(t, "response body", string(decompressedData))
	})
}
