package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

type fakeTaskStore struct {
	createdInput tasks.CreateInput
	listFilter   tasks.ListFilter
	task         tasks.Task
	items        []tasks.Task
	err          error
}

func (store *fakeTaskStore) Create(_ context.Context, input tasks.CreateInput) (tasks.Task, error) {
	store.createdInput = input
	return store.task, store.err
}

func (store *fakeTaskStore) Get(context.Context, int64) (tasks.Task, error) {
	return store.task, store.err
}

func (store *fakeTaskStore) List(_ context.Context, filter tasks.ListFilter) ([]tasks.Task, error) {
	store.listFilter = filter
	return store.items, store.err
}

func (store *fakeTaskStore) Update(context.Context, int64, tasks.UpdateInput, scoring.Result) (tasks.Task, error) {
	return store.task, store.err
}

func (store *fakeTaskStore) Confirm(context.Context, int64, scoring.Result) (tasks.Task, error) {
	return store.task, store.err
}

func (store *fakeTaskStore) Publish(context.Context, int64) (tasks.Task, error) {
	return store.task, store.err
}

func TestCreateTask(t *testing.T) {
	store := &fakeTaskStore{task: tasks.Task{ID: 42, Status: "draft"}}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"  Need a dashboard  "}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

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

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

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

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestCreateTaskRejectsUnknownFields(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"test","status":"published"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestCreateTaskRejectsMultipleJSONObjects(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"initial_description":"first"}{"initial_description":"second"}`))
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestGetTaskInternalError(t *testing.T) {
	store := &fakeTaskStore{err: errors.New("database error")}
	request := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestConfirmTaskRejectsPublishedTask(t *testing.T) {
	store := &fakeTaskStore{task: tasks.Task{ID: 1, Status: "published"}}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks/1/confirm", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
}

func TestPublishTaskRequiresConfirmation(t *testing.T) {
	store := &fakeTaskStore{task: tasks.Task{ID: 1, Status: "draft"}}
	request := httptest.NewRequest(http.MethodPost, "/api/tasks/1/publish", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
}

func TestListTasksPassesCatalogFilters(t *testing.T) {
	store := &fakeTaskStore{items: []tasks.Task{}}
	request := httptest.NewRequest(http.MethodGet, "/api/tasks?topic=analytics&readiness_level=ready&sort=readiness_asc", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if store.listFilter.Topic != "analytics" || store.listFilter.ReadinessLevel != "ready" || store.listFilter.Sort != "readiness_asc" {
		t.Fatalf("unexpected filters: %+v", store.listFilter)
	}
	if recorder.Body.String() != "[]\n" {
		t.Fatalf("expected empty JSON array, got %s", recorder.Body.String())
	}
}

func TestListTasksRejectsInvalidFilter(t *testing.T) {
	store := &fakeTaskStore{}
	request := httptest.NewRequest(http.MethodGet, "/api/tasks?readiness_level=excellent", nil)
	recorder := httptest.NewRecorder()

	New(fakeDatabase{}, store, nil, nil).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
