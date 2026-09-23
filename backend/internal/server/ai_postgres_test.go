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
	"github.com/shmaloogles/business-task-platform/backend/internal/proposals"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/tasks"
)

// Opt-in integration test. Requires the existing migration; deletes only its own fixtures.
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
	// Own fixtures: this test does not depend on seed IDs or change demo records.
	var teamIDs []int64
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.Exec(cleanup, "DELETE FROM teams WHERE id = ANY($1)", teamIDs); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i < 3; i++ {
		var id int64
		if err := db.QueryRow(ctx, "INSERT INTO teams (name) VALUES ($1) RETURNING id", fmt.Sprintf("integration team %d", i)).Scan(&id); err != nil {
			t.Fatal(err)
		}
		teamIDs = append(teamIDs, id)
	}
	mux := New(db, tasks.NewStore(db), nil, proposals.NewStore(db))
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
	var preview scoring.Result
	request("POST", "/api/tasks/score", generated.Card, &preview, http.StatusOK)
	if !reflect.DeepEqual(preview, generated.Readiness) {
		t.Fatal("preview and generation readiness differ")
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
	request("POST", "/api/tasks/score", generated.Card, &preview, http.StatusOK)
	request("GET", path, nil, &loaded, http.StatusOK)
	if preview.Score != 50 || loaded.ReadinessScore != 60 || loaded.Users == nil {
		t.Fatal("preview must reflect edits without persisting them")
	}
	request("PUT", path, generated.Card, &updated, http.StatusOK)
	request("GET", path, nil, &loaded, http.StatusOK)
	if loaded.ReadinessScore != int16(preview.Score) {
		t.Fatal("persisted readiness differs from preview")
	}
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

	// Decisions are independent: accepting a second team must preserve the first.
	wantStatuses := []string{"accepted", "accepted", "rejected"}
	for i, teamID := range teamIDs {
		var proposal proposals.Proposal
		request("POST", path+"/proposals", proposals.CreateInput{
			TeamID: teamID, SolutionIdea: "Test solution", Plan: "Build and verify", Timeline: "2 weeks",
		}, &proposal, http.StatusCreated)
		if proposal.Status != "pending" || proposal.TaskID != created.ID || proposal.TeamID != teamID {
			t.Fatalf("unexpected new proposal: %+v", proposal)
		}
		request("PATCH", fmt.Sprintf("/api/proposals/%d", proposal.ID), proposals.StatusInput{Status: wantStatuses[i]}, &proposal, http.StatusOK)
		if proposal.Status != wantStatuses[i] {
			t.Fatalf("decision not applied: %+v", proposal)
		}
	}
	var listed []proposals.Proposal
	request("GET", path+"/proposals", nil, &listed, http.StatusOK)
	statuses := map[int64]string{}
	for _, item := range listed {
		statuses[item.TeamID] = item.Status
	}
	if len(listed) != 3 {
		t.Fatalf("expected 3 proposals, got %d", len(listed))
	}
	for i, id := range teamIDs {
		if statuses[id] != wantStatuses[i] {
			t.Fatalf("persisted decision mismatch for team %d", id)
		}
	}
	var allProposals []proposals.Proposal
	request("GET", "/api/proposals", nil, &allProposals, http.StatusOK)
	for _, expected := range listed {
		found := false
		for _, actual := range allProposals {
			if actual.ID == expected.ID {
				found = true
				if !reflect.DeepEqual(actual, expected) {
					t.Fatal("global proposal list differs from task list")
				}
			}
		}
		if !found {
			t.Fatalf("proposal %d missing from global list", expected.ID)
		}
	}

	// Zero readiness does not prevent publication or proposals.
	var empty tasks.Task
	request("POST", "/api/tasks", tasks.CreateInput{InitialDescription: "Incomplete test description"}, &empty, http.StatusCreated)
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := db.Exec(cleanup, "DELETE FROM tasks WHERE id = $1", empty.ID); err != nil {
			t.Error(err)
		}
	}()
	emptyPath := fmt.Sprintf("/api/tasks/%d", empty.ID)
	// The business list includes drafts; the public catalog must not expose them.
	request("GET", "/api/tasks?scope=all", nil, &catalog, http.StatusOK)
	found = false
	for _, task := range catalog {
		if task.ID == empty.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("draft missing from business list")
	}
	request("GET", "/api/tasks", nil, &catalog, http.StatusOK)
	for _, task := range catalog {
		if task.ID == empty.ID {
			t.Fatal("draft leaked into public catalog")
		}
	}
	// Invalid JSON must not clear a saved card.
	r, err := http.NewRequestWithContext(ctx, "PUT", server.URL+path, bytes.NewBufferString("null"))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("null PUT returned %d", resp.StatusCode)
	}
	request("GET", path, nil, &loaded, http.StatusOK)
	if !reflect.DeepEqual(loaded.Card(), generated.Card) || loaded.ReadinessScore != 50 {
		t.Fatal("invalid PUT changed stored card")
	}
	input := proposals.CreateInput{TeamID: teamIDs[0], SolutionIdea: "Test", Plan: "Test", Timeline: "1 week"}
	request("POST", emptyPath+"/proposals", input, &rejected, http.StatusConflict)
	if rejected.Error.Code != "task_not_published" {
		t.Fatalf("unexpected error: %+v", rejected)
	}
	request("POST", emptyPath+"/confirm", nil, &loaded, http.StatusOK)
	request("POST", emptyPath+"/publish", nil, &loaded, http.StatusOK)
	if loaded.ReadinessScore != 0 || loaded.Status != "published" {
		t.Fatalf("zero-score publication failed: %+v", loaded)
	}
	request("GET", "/api/tasks?readiness_level=draft", nil, &catalog, http.StatusOK)
	found = false
	for _, task := range catalog {
		if task.ID == empty.ID {
			found = true
		}
	}
	if !found {
		t.Fatal("published zero-score task missing from catalog")
	}
	var proposal proposals.Proposal
	request("POST", emptyPath+"/proposals", input, &proposal, http.StatusCreated)
	if proposal.Status != "pending" {
		t.Fatalf("zero-score proposal failed: %+v", proposal)
	}
}
