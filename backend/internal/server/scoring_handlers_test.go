package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/ai"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
)

func TestScorePreview(t *testing.T) {
	// No store or provider is needed to score an edited card.
	mux := New(fakeDatabase{}, nil, nil)
	RegisterAIRoutes(mux, ai.New(aiTestProvider(func(context.Context) ([]byte, error) {
		t.Fatal("scoring must not call AI")
		return nil, nil
	}), 0, false))
	for _, tc := range []struct {
		name, input string
		score       int
		level       string
	}{
		{"empty", `{}`, 0, "Draft"},
		{"partial", `{"context":"Процесс","contact":"team@example.test"}`, 15, "Draft"},
		{"edited", `{"context":"Процесс","need":"Изменение","data":"CSV","contact":"team@example.test"}`, 45, "Workable"},
		{"cleared", `{"context":null,"contact":" \t\n"}`, 0, "Draft"},
		{"full", `{"context":"x","need":"x","data":"x","users":"x","constraints":"x","expected_result":"x","success_criteria":"x","contact":"x","interaction_format":"x"}`, 100, "Priority"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			mux.ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/tasks/score", strings.NewReader(tc.input)))
			var result scoring.Result
			if r.Code != http.StatusOK || json.Unmarshal(r.Body.Bytes(), &result) != nil {
				t.Fatalf("%d %s", r.Code, r.Body)
			}
			if result.Score != tc.score || result.Level != tc.level || len(result.Breakdown) != 7 || len(result.MissingFields) != len(result.Suggestions) {
				t.Fatalf("unexpected preview: %+v", result)
			}
		})
	}
}

func TestScorePreviewRejectsInvalidInput(t *testing.T) {
	mux := New(fakeDatabase{}, nil, nil)
	RegisterAIRoutes(mux, nil)
	for _, body := range []string{
		`null`, `[]`, `{`, `{} {}`, `{"users":[]}`, `{"data":42}`,
		`{"score":100}`, `{"status":"confirmed"}`, `{"card":{}}`,
	} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/tasks/score", strings.NewReader(body)))
		if r.Code != http.StatusBadRequest {
			t.Fatalf("body %s: %d %s", body, r.Code, r.Body)
		}
	}
	r := httptest.NewRecorder()
	body := `{"context":"` + strings.Repeat("x", 1<<20) + `"}`
	mux.ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/tasks/score", strings.NewReader(body)))
	if r.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("%d %s", r.Code, r.Body)
	}
}
