package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to read gzip data", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = gzipReader
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gzipWriter := gzip.NewWriter(w)
			w.Header().Set("Content-Encoding", "gzip")
			defer gzipWriter.Close()
			gzipResponseWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			next.ServeHTTP(gzipResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}
