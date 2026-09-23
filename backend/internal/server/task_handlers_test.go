package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

type fakeTaskStore struct {
	createdInput tasks.CreateInput
	task         tasks.Task
	err          error
}

func (store *fakeTaskStore) Create(_ context.Context, input tasks.CreateInput) (tasks.Task, error) {
	store.createdInput = input
	return store.task, store.err
}

func (store *fakeTaskStore) Get(context.Context, int64) (tasks.Task, error) {
	return store.task, store.err
}

func (store *fakeTaskStore) Update(context.Context, int64, tasks.UpdateInput) (tasks.Task, error) {
	return store.task, store.err
}

func TestCreateTask(t *testing.T) {
	store := &fakeTaskStore{task: tasks.Task{ID: 42, Status: "draft"}}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"  Need a dashboard  "}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}
	if store.createdInput.InitialDescription != "Need a dashboard" {
		t.Fatalf("unexpected description: %q", store.createdInput.InitialDescription)
	}
}

func TestCreateTaskRequiresDescription(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"  "}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"description_required"`) {
		t.Fatalf("unexpected response body: %s", recorder.Body.String())
	}
}

func TestGetTaskNotFound(t *testing.T) {
	store := &fakeTaskStore{err: tasks.ErrNotFound}
	request := httptest.NewRequest(http.MethodGet, "/api/tasks/999", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestCreateTaskRejectsUnknownFields(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"test","status":"published"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestCreateTaskRejectsMultipleJSONObjects(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"first"}{"initial_description":"second"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestGetTaskInternalError(t *testing.T) {
	store := &fakeTaskStore{err: errors.New("database error")}
	request := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}
