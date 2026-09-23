package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type response struct {
	Status string `json:"status"`
}

type healthChecker interface {
	Ping(context.Context) error
}

func New(database healthChecker, tasks taskStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api", apiHandler)
	mux.HandleFunc("GET /health", healthHandler(database))
	mux.HandleFunc("POST /api/tasks", createTaskHandler(tasks))
	mux.HandleFunc("GET /api/tasks/{id}", getTaskHandler(tasks))
	mux.HandleFunc("PUT /api/tasks/{id}", updateTaskHandler(tasks))

	return mux
}

func apiHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, response{Status: "ok"})
}

func healthHandler(database healthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, response{Status: "unhealthy"})
			return
		}

		writeJSON(w, http.StatusOK, response{Status: "healthy"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
