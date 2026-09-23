package server

import (
	"context"
	"net/http"

	"github.com/shmaloogles/business-task-platform/backend/internal/teams"
)

type teamStore interface {
	List(context.Context) ([]teams.Team, error)
}

func listTeamsHandler(store teamStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := store.List(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not list teams")
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}
