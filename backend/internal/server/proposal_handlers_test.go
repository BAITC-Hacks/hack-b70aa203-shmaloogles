package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/proposals"
)

type fakeProposalStore struct {
	proposal proposals.Proposal
	items    []proposals.Proposal
	err      error
}

func (store fakeProposalStore) Create(context.Context, int64, proposals.CreateInput) (proposals.Proposal, error) {
	return store.proposal, store.err
}

func (store fakeProposalStore) ListByTask(context.Context, int64) ([]proposals.Proposal, error) {
	return store.items, store.err
}

func (store fakeProposalStore) UpdateStatus(context.Context, int64, string) (proposals.Proposal, error) {
	return store.proposal, store.err
}

func TestCreateProposal(t *testing.T) {
	store := fakeProposalStore{proposal: proposals.Proposal{ID: 10, TaskID: 1, TeamID: 2, Status: "pending"}}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks/1/proposals", strings.NewReader(`{"team_id":2,"solution_idea":"Dashboard","plan":"Build and test","timeline":"3 weeks"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, nil, nil, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}
}

func TestCreateProposalRequiresFields(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/tasks/1/proposals", strings.NewReader(`{"team_id":2}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, nil, nil, fakeProposalStore{}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestUpdateProposalRejectsPendingStatus(t *testing.T) {
	request := httptest.NewRequest(http.MethodPatch, "/api/proposals/1", strings.NewReader(`{"status":"pending"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, nil, nil, fakeProposalStore{}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
