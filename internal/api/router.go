package api

import (
	"database/sql"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Role, X-User-ID")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func NewRouter(db *sql.DB) http.Handler {
	handlers := NewHandlers(db)
	
	mux := http.NewServeMux()
	
	mux.Handle("/prescriptions", corsMiddleware(RoleMiddleware(http.HandlerFunc(handlers.CreatePrescription))))
	mux.Handle("/analytics/top-drugs", corsMiddleware(RoleMiddleware(http.HandlerFunc(handlers.GetTopDrugs))))
	
	return mux
}