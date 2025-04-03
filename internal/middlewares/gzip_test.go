package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_CompressResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	})
	handlerWithMiddleware := GzipMiddleware(handler)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
	gzipReader, err := gzip.NewReader(rr.Body)
	assert.NoError(t, err)
	defer gzipReader.Close()
	decompressedBody := new(bytes.Buffer)
	_, err = decompressedBody.ReadFrom(gzipReader)
	assert.NoError(t, err)
	assert.Equal(t, "", decompressedBody.String())
}

func TestGzipMiddleware_DecompressRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Write(body)
	})
	handlerWithMiddleware := GzipMiddleware(handler)
	var originalBody = "This is a gzipped body"
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte(originalBody))
	assert.NoError(t, err)
	err = gzipWriter.Close()
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	req.ContentLength = int64(buf.Len())
	rr := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(rr, req)
	assert.Equal(t, originalBody, rr.Body.String())
}

func TestGzipMiddleware_NoCompressionWhenNotNeeded(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	})
	handlerWithMiddleware := GzipMiddleware(handler)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(rr, req)
	assert.Empty(t, rr.Header().Get("Content-Encoding"))
	assert.Equal(t, "", rr.Body.String())
}
