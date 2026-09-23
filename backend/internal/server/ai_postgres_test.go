package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/ai"
	"github.com/shmaloogles/business-task-platform/backend/internal/database"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

// Opt-in integration test. Requires the existing migration; deletes only its own row.
func TestAIFlowPostgres(t *testing.T) {
	dsn := os.Getenv("AI_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set AI_TEST_DATABASE_URL to run against a migrated test database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := database.Open(ctx, dsn)
	if err != nil {
		t.Fatal("test database connection failed")
	}
	defer db.Close()
	mux := New(db, tasks.NewStore(db), nil, nil)
	RegisterAIRoutes(mux, ai.New(nil, 0, false))
	server := httptest.NewServer(mux)
	defer server.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	request := func(method, path string, input, output any, status int) {
		t.Helper()
		var body io.Reader
		if input != nil {
			raw, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			body = bytes.NewReader(raw)
		}
		r, err := http.NewRequestWithContext(ctx, method, server.URL+path, body)
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != status {
			t.Fatalf("%s %s: status %d, want %d", method, path, resp.StatusCode, status)
		}
		if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
			t.Fatal(err)
		}
	}
	description := "Заказы кафе ведём вручную в таблице."
	var created tasks.Task
	request("POST", "/api/tasks", tasks.CreateInput{InitialDescription: description}, &created, http.StatusCreated)
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.Exec(cleanup, "DELETE FROM tasks WHERE id = $1", created.ID); err != nil {
			t.Error("could not remove test task")
		}
	}()
	var questions ai.Clarification
	request("POST", "/api/tasks/clarify", map[string]string{"description": description}, &questions, http.StatusOK)
	if questions.Mode != "mock" || len(questions.Questions) < 3 {
		t.Fatalf("%+v", questions)
	}
	var generated struct {
		ai.Generation
		Readiness scoring.Result `json:"readiness"`
	}
	request("POST", "/api/tasks/generate", ai.GenerateInput{Description: description, Answers: []ai.Answer{
		{Field: "need", Answer: "Сократить ручной ввод"},
		{Field: "users", Answer: "Администраторы кафе"},
		{Field: "expected_result", Answer: "Прототип формы заказа"},
		{Field: "success_criteria", Answer: "Оформление заказа за 1 минуту"},
	}}, &generated, http.StatusOK)
	if generated.Mode != "mock" || generated.Readiness.Score != 60 {
		t.Fatalf("%+v", generated)
	}
	path := fmt.Sprintf("/api/tasks/%d", created.ID)
	var updated, loaded tasks.Task
	request("PUT", path, generated.Card, &updated, http.StatusOK)
	request("GET", path, nil, &loaded, http.StatusOK)
	if !reflect.DeepEqual(loaded.Card(), generated.Card) || !reflect.DeepEqual(updated.Card(), loaded.Card()) {
		t.Fatal("card changed across PostgreSQL round trip")
	}
	if loaded.InitialDescription == nil || *loaded.InitialDescription != description || loaded.Status != "draft" || loaded.ConfirmedAt != nil || loaded.ReadinessScore != 60 || loaded.ReadinessLevel != "workable" {
		t.Fatalf("unexpected lifecycle state: %+v", loaded)
	}
	loaded.RecalculateReadiness()
	if loaded.ReadinessScore != 60 || loaded.ReadinessLevel != "workable" {
		t.Fatalf("%+v", loaded)
	}
	// Clear a field through the real PUT to verify null handling and a lower score.
	generated.Card.Users = nil
	request("PUT", path, generated.Card, &updated, http.StatusOK)
	request("GET", path, nil, &loaded, http.StatusOK)
	loaded.RecalculateReadiness()
	if loaded.Users != nil || loaded.ReadinessScore != 50 {
		t.Fatalf("null/edit did not survive persistence: %+v", loaded)
	}
	var rejected errorResponse
	request("POST", path+"/publish", nil, &rejected, http.StatusConflict)
	request("POST", path+"/confirm", nil, &loaded, http.StatusOK)
	if loaded.Status != "confirmed" || loaded.ConfirmedAt == nil || loaded.ReadinessScore != 50 {
		t.Fatalf("confirmation failed: %+v", loaded)
	}
	request("POST", path+"/publish", nil, &loaded, http.StatusOK)
	if loaded.Status != "published" || loaded.PublishedAt == nil {
		t.Fatalf("publication failed: %+v", loaded)
	}
	var catalog []tasks.Task
	request("GET", "/api/tasks?readiness_level=workable&sort=readiness_desc", nil, &catalog, http.StatusOK)
	found := false
	for _, task := range catalog {
		if task.ID == created.ID {
			found = true
			if task.ReadinessScore != 50 {
				t.Fatal("catalog score mismatch")
			}
		}
	}
	if !found {
		t.Fatal("published AI task missing from catalog")
	}
}
