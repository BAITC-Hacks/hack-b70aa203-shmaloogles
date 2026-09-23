package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/teams"
)

type fakeTeamStore struct {
	items []teams.Team
	err   error
}

func (store fakeTeamStore) List(context.Context) ([]teams.Team, error) {
	return store.items, store.err
}

func TestListTeams(t *testing.T) {
	store := fakeTeamStore{items: []teams.Team{{ID: 1, Name: "Data Sparks"}}}
	request := httptest.NewRequest(http.MethodGet, "/api/teams", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, nil, store, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if got := recorder.Body.String(); got == "[]\n" {
		t.Fatalf("expected team in response, got %s", got)
	}
}
