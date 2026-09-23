package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/proposals"
)

type proposalStore interface {
	Create(context.Context, int64, proposals.CreateInput) (proposals.Proposal, error)
	List(context.Context) ([]proposals.Proposal, error)
	ListByTask(context.Context, int64) ([]proposals.Proposal, error)
	UpdateStatus(context.Context, int64, string) (proposals.Proposal, error)
}

func listAllProposalsHandler(store proposalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := store.List(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not list proposals")
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func createProposalHandler(store proposalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, ok := taskID(w, r)
		if !ok {
			return
		}
		var input proposals.CreateInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
			return
		}
		input.SolutionIdea = strings.TrimSpace(input.SolutionIdea)
		input.Plan = strings.TrimSpace(input.Plan)
		input.Timeline = strings.TrimSpace(input.Timeline)
		if input.TeamID < 1 || input.SolutionIdea == "" || input.Plan == "" || input.Timeline == "" {
			writeError(w, http.StatusBadRequest, "invalid_proposal", "team_id, solution_idea, plan, and timeline are required")
			return
		}
		if input.PrototypeURL != nil {
			trimmed := strings.TrimSpace(*input.PrototypeURL)
			if trimmed == "" {
				input.PrototypeURL = nil
			} else {
				input.PrototypeURL = &trimmed
			}
		}

		proposal, err := store.Create(r.Context(), taskID, input)
		switch {
		case errors.Is(err, proposals.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, "task_not_found", "task not found")
		case errors.Is(err, proposals.ErrTaskUnavailable):
			writeError(w, http.StatusConflict, "task_not_published", "proposals are allowed only for published tasks")
		case errors.Is(err, proposals.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, "team_not_found", "team not found")
		case err != nil:
			writeError(w, http.StatusInternalServerError, "internal_error", "could not create proposal")
		default:
			writeJSON(w, http.StatusCreated, proposal)
		}
	}
}

func listProposalsHandler(store proposalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, ok := taskID(w, r)
		if !ok {
			return
		}
		items, err := store.ListByTask(r.Context(), taskID)
		if errors.Is(err, proposals.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task_not_found", "task not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not list proposals")
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func updateProposalHandler(store proposalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil || id < 1 {
			writeError(w, http.StatusBadRequest, "invalid_proposal_id", "proposal id must be a positive integer")
			return
		}
		var input proposals.StatusInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
			return
		}
		if input.Status != "accepted" && input.Status != "rejected" {
			writeError(w, http.StatusBadRequest, "invalid_proposal_status", "status must be accepted or rejected")
			return
		}

		proposal, err := store.UpdateStatus(r.Context(), id, input.Status)
		if errors.Is(err, proposals.ErrNotFound) {
			writeError(w, http.StatusNotFound, "proposal_not_found", "proposal not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not update proposal")
			return
		}
		writeJSON(w, http.StatusOK, proposal)
	}
}
