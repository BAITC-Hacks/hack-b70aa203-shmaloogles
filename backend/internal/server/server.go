package server

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status string `json:"status"`
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api", apiHandler)
	mux.HandleFunc("GET /health", healthHandler)

	return mux
}

func apiHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, response{Status: "ok"})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, response{Status: "healthy"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
