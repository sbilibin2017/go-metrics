package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type Config interface {
	GetKey() string
}

func HMACMiddleware(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secretKey := cfg.GetKey()
			if secretKey == "" {
				next.ServeHTTP(w, r)
				return
			}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read request body", http.StatusInternalServerError)
				return
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			hashHeader := r.Header.Get("HashSHA256")
			expectedHash := computeHMAC(bodyBytes, secretKey)
			if hashHeader == "" || hashHeader != expectedHash {
				http.Error(w, "Invalid hash signature", http.StatusBadRequest)
				return
			}
			recorder := &responseRecorder{ResponseWriter: w, bodyBuffer: &bytes.Buffer{}}
			next.ServeHTTP(recorder, r)
			responseHash := computeHMAC(recorder.bodyBuffer.Bytes(), secretKey)
			w.Header().Set("HashSHA256", responseHash)
			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.statusCode)
			w.Write(recorder.bodyBuffer.Bytes())
		})
	}
}

func computeHMAC(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

type responseRecorder struct {
	http.ResponseWriter
	bodyBuffer *bytes.Buffer
	statusCode int
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	return rw.bodyBuffer.Write(b)
}

func (rw *responseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}
