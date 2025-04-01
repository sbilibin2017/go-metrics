package handlers

import (
	"go-metrics/internal/engines"
	"net/http"
)

func PingDBHandler(db *engines.DBEngine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
	}
}
