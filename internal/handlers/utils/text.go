package utils

import "net/http"

func SendTextResponse(w http.ResponseWriter, statusCode int, message []byte) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write(message)
}
