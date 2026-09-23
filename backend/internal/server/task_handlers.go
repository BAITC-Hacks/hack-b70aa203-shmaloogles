package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

type taskStore interface {
	Create(context.Context, tasks.CreateInput) (tasks.Task, error)
	Get(context.Context, int64) (tasks.Task, error)
	List(context.Context, tasks.ListFilter) ([]tasks.Task, error)
	Update(context.Context, int64, tasks.UpdateInput, scoring.Result) (tasks.Task, error)
	Confirm(context.Context, int64, scoring.Result) (tasks.Task, error)
	Publish(context.Context, int64) (tasks.Task, error)
}

func listTasksHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := tasks.ListFilter{
			Topic:          strings.TrimSpace(r.URL.Query().Get("topic")),
			ReadinessLevel: strings.TrimSpace(r.URL.Query().Get("readiness_level")),
			Sort:           strings.TrimSpace(r.URL.Query().Get("sort")),
			IncludeAll:     strings.TrimSpace(r.URL.Query().Get("scope")) == "all",
		}
		if scope := strings.TrimSpace(r.URL.Query().Get("scope")); scope != "" && scope != "all" {
			writeError(w, http.StatusBadRequest, "invalid_scope", "scope must be all when provided")
			return
		}
		if !validReadinessLevel(filter.ReadinessLevel) {
			writeError(w, http.StatusBadRequest, "invalid_readiness_level", "readiness_level must be draft, workable, ready, or priority")
			return
		}
		if filter.Sort != "" && filter.Sort != "readiness_desc" && filter.Sort != "readiness_asc" {
			writeError(w, http.StatusBadRequest, "invalid_sort", "sort must be readiness_desc or readiness_asc")
			return
		}

		items, err := store.List(r.Context(), filter)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not list tasks")
			return
		}
		writeJSON(w, http.StatusOK, items)
	}
}

func validReadinessLevel(level string) bool {
	switch level {
	case "", "draft", "workable", "ready", "priority":
		return true
	default:
		return false
	}
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

		readiness := scoring.Calculate(input)
		task, err := store.Update(r.Context(), id, input, readiness)
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

func confirmTaskHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}

		task, err := store.Get(r.Context(), id)
		if !handleTaskLookupError(w, err) {
			return
		}
		if task.Status != "draft" {
			writeError(w, http.StatusConflict, "invalid_task_status", "only a draft task can be confirmed")
			return
		}

		task, err = store.Confirm(r.Context(), id, scoring.Calculate(task.Card()))
		if errors.Is(err, tasks.ErrNotFound) {
			writeError(w, http.StatusConflict, "invalid_task_status", "task status changed; reload and try again")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not confirm task")
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func publishTaskHandler(store taskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}

		task, err := store.Get(r.Context(), id)
		if !handleTaskLookupError(w, err) {
			return
		}
		if task.Status != "confirmed" {
			writeError(w, http.StatusConflict, "invalid_task_status", "only a confirmed task can be published")
			return
		}

		task, err = store.Publish(r.Context(), id)
		if errors.Is(err, tasks.ErrNotFound) {
			writeError(w, http.StatusConflict, "invalid_task_status", "task status changed; reload and try again")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "could not publish task")
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

func handleTaskLookupError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, tasks.ErrNotFound) {
		writeError(w, http.StatusNotFound, "task_not_found", "task not found")
		return false
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "could not get task")
		return false
	}
	return true
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
