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
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(io.Reader(bytes.NewReader(bodyBytes)))
			key := cfg.GetKey()
			hash := computeHMAC(bodyBytes, key)
			r.Header.Set("HashSHA256", hash)
			next.ServeHTTP(w, r)
		})
	}
}

func computeHMAC(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
