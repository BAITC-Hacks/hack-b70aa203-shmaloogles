package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/ai"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

type aiTestProvider func(context.Context) ([]byte, error)

func (p aiTestProvider) Complete(ctx context.Context, _ string, _, _ json.RawMessage) ([]byte, error) {
	return p(ctx)
}

type recordingAIStore struct {
	fakeTaskStore
	input tasks.UpdateInput
}

func (s *recordingAIStore) Update(_ context.Context, _ int64, input tasks.UpdateInput) (tasks.Task, error) {
	s.input = input
	return s.task, nil
}

func TestAIHTTPToTaskUpdate(t *testing.T) {
	store := &recordingAIStore{fakeTaskStore: fakeTaskStore{task: tasks.Task{ID: 42, Status: "draft"}}}
	mux := New(fakeDatabase{}, store)
	RegisterAIRoutes(mux, ai.New(nil, 0, false))
	q := httptest.NewRecorder()
	mux.ServeHTTP(q, httptest.NewRequest("POST", "/api/tasks/clarify", strings.NewReader(`{"description":"Заказы в таблице"}`)))
	var questions ai.Clarification
	if err := json.Unmarshal(q.Body.Bytes(), &questions); err != nil || q.Code != 200 || questions.Mode != "mock" || len(questions.Questions) < 3 {
		t.Fatalf("%d %s %v", q.Code, q.Body, err)
	}
	r := httptest.NewRecorder()
	mux.ServeHTTP(r, httptest.NewRequest("POST", "/api/tasks/generate", strings.NewReader(`{"description":"Заказы в таблице","answers":[{"field":"need","answer":"Упростить ввод"}]}`)))
	var generated struct {
		ai.Generation
		Readiness scoring.Result `json:"readiness"`
	}
	if err := json.Unmarshal(r.Body.Bytes(), &generated); err != nil || r.Code != 200 {
		t.Fatalf("%d %s %v", r.Code, r.Body, err)
	}
	if generated.Mode != "mock" || generated.Readiness.Score != 20 || generated.Card.Context == nil || generated.Card.Need == nil || generated.Card.Contact != nil {
		t.Fatalf("%+v", generated)
	}
	raw, _ := json.Marshal(generated.Card)
	updated := httptest.NewRecorder()
	mux.ServeHTTP(updated, httptest.NewRequest("PUT", "/api/tasks/42", bytes.NewReader(raw)))
	if updated.Code != 200 || !reflect.DeepEqual(store.input, generated.Card) {
		t.Fatalf("card not accepted by CRUD: %d %s", updated.Code, updated.Body)
	}
	if store.task.Status != "draft" || store.task.ReadinessScore != 0 {
		t.Fatal("AI must not confirm or persist a score")
	}
}

func TestAIHTTPInvalidInput(t *testing.T) {
	for _, path := range []string{"/api/tasks/clarify", "/api/tasks/generate"} {
		for _, tc := range []struct {
			name, body string
			status     int
		}{
			{"broken", "{", 400}, {"null", "null", 400}, {"empty", `{}`, 400},
			{"wrong type", `{"description":42}`, 400},
			{"unknown", `{"description":"x","status":"published"}`, 400},
			{"trailing", `{"description":"x"}{}`, 400},
			{"too large", `{"description":"` + strings.Repeat("x", (1<<20)+1) + `"}`, 413},
		} {
			t.Run(path+tc.name, func(t *testing.T) {
				mux := New(fakeDatabase{}, nil)
				RegisterAIRoutes(mux, ai.New(nil, 0, false))
				r := httptest.NewRecorder()
				mux.ServeHTTP(r, httptest.NewRequest("POST", path, strings.NewReader(tc.body)))
				var body errorResponse
				if json.Unmarshal(r.Body.Bytes(), &body) != nil || r.Code != tc.status || body.Error.Code == "" {
					t.Fatalf("%d %s", r.Code, r.Body)
				}
			})
		}
	}
}

func TestAIHTTPErrorsAndFallback(t *testing.T) {
	for _, tc := range []struct {
		name     string
		provider ai.Provider
		status   int
		code     string
	}{
		{"provider", aiTestProvider(func(context.Context) ([]byte, error) { return nil, errors.New("private upstream body") }), 502, "ai_unavailable"},
		{"invalid", aiTestProvider(func(context.Context) ([]byte, error) { return []byte(`{}`), nil }), 502, "invalid_ai_response"},
		{"timeout", aiTestProvider(func(ctx context.Context) ([]byte, error) { <-ctx.Done(); return nil, ctx.Err() }), 504, "ai_timeout"},
	} {
		for _, path := range []string{"/api/tasks/clarify", "/api/tasks/generate"} {
			for _, fallback := range []bool{false, true} {
				mux := New(fakeDatabase{}, nil)
				RegisterAIRoutes(mux, ai.New(tc.provider, time.Millisecond, fallback))
				r := httptest.NewRecorder()
				mux.ServeHTTP(r, httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"description":"x"}`)))
				if fallback {
					var metadata ai.Metadata
					if json.Unmarshal(r.Body.Bytes(), &metadata) != nil || r.Code != 200 || metadata.Mode != "fallback" || metadata.FallbackReason == "" {
						t.Fatalf("%d %s", r.Code, r.Body)
					}
				} else {
					var body errorResponse
					if json.Unmarshal(r.Body.Bytes(), &body) != nil || r.Code != tc.status || body.Error.Code != tc.code {
						t.Fatalf("%d %s", r.Code, r.Body)
					}
				}
				if strings.Contains(r.Body.String(), "private") {
					t.Fatal("leaked upstream error")
				}
			}
		}
	}
}
