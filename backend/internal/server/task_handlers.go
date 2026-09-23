package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

type taskStore interface {
	Create(context.Context, tasks.CreateInput) (tasks.Task, error)
	Get(context.Context, int64) (tasks.Task, error)
	Update(context.Context, int64, tasks.UpdateInput) (tasks.Task, error)
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func createTaskHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input tasks.CreateInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
			return
		}

		input.InitialDescription = strings.TrimSpace(input.InitialDescription)
		if input.InitialDescription == "" {
			writeError(w, http.StatusBadRequest, "description_required", "initial_description is required")
			return
		}

		task, err := store.Create(r.Context(), input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not create task")
			return
		}
		writeJSON(w, http.StatusCreated, task)
	}
}

func getTaskHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}

		task, err := store.Get(r.Context(), id)
		if errors.Is(err, tasks.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task_not_found", "task not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not get task")
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func updateTaskHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}

		var input tasks.UpdateInput
		if err := decodeJSON(r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
			return
		}

		task, err := store.Update(r.Context(), id, input)
		if errors.Is(err, tasks.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task_not_found", "task not found")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not update task")
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func taskID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid_task_id", "task id must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeJSON(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}
