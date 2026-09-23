package tasks

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

func TestCardAndUpdateContract(t *testing.T) {
	fields := map[string]string{}
	for _, field := range taskcard.Fields() {
		fields[field] = "value for " + field
	}
	raw, _ := json.Marshal(fields)
	var task Task
	if err := json.Unmarshal(raw, &task); err != nil {
		t.Fatal(err)
	}
	card := task.Card()
	for field, value := range card.Values() {
		if value == nil || *value != fields[field] {
			t.Fatalf("lost %s", field)
		}
	}
	// Generated cards are accepted directly by Store.Update without conversion.
	var update UpdateInput = card
	encoded, err := json.Marshal(update)
	if err != nil || string(encoded) == "" {
		t.Fatal(err)
	}
	var roundTrip map[string]string
	if err := json.Unmarshal(encoded, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fields, roundTrip) {
		t.Fatalf("contract changed: %s", encoded)
	}
}

func TestRecalculateUsesDatabaseContract(t *testing.T) {
	now := time.Now()
	text := "provided"
	task := Task{ID: 42, Status: "confirmed", ConfirmedAt: &now, CreatedAt: now, UpdatedAt: now,
		Context: &text, Need: &text, Users: &text, Data: &text, Constraints: &text,
		ExpectedResult: &text, SuccessCriteria: &text, Contact: &text, InteractionFormat: &text,
	}
	task.RecalculateReadiness()
	if task.ReadinessScore != 100 || task.ReadinessLevel != "priority" {
		t.Fatalf("%+v", task)
	}
	var breakdown []scoring.Category
	if err := json.Unmarshal(task.ReadinessBreakdown, &breakdown); err != nil || len(breakdown) != 7 {
		t.Fatalf("%s %v", task.ReadinessBreakdown, err)
	}
	if !reflect.DeepEqual(task.MissingInformation, []string{"title", "topic"}) || len(task.Suggestions) != 2 {
		t.Fatalf("%+v", task)
	}
	if task.ID != 42 || task.Status != "confirmed" || task.ConfirmedAt != &now || !task.UpdatedAt.Equal(now) || !task.CreatedAt.Equal(now) || task.PublishedAt != nil {
		t.Fatal("changed lifecycle metadata")
	}
	task.Data = nil
	task.RecalculateReadiness()
	if task.ReadinessScore != 80 || task.ReadinessLevel != "ready" || len(task.Suggestions) != 3 {
		t.Fatalf("stale score: %+v", task)
	}
	task.ExpectedResult = nil
	task.RecalculateReadiness()
	if task.ReadinessScore != 65 || task.ReadinessLevel != "workable" {
		t.Fatalf("%+v", task)
	}
	empty := Task{Status: "draft"}
	empty.RecalculateReadiness()
	if empty.ReadinessScore != 0 || empty.ReadinessLevel != "draft" || empty.Status != "draft" || empty.ConfirmedAt != nil || len(empty.MissingInformation) != 11 {
		t.Fatalf("%+v", empty)
	}
}
