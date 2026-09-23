package tasks

import (
	"encoding/json"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

func (task Task) Card() taskcard.Card {
	return taskcard.Card{
		Title: task.Title, Topic: task.Topic, Context: task.Context, Need: task.Need,
		Users: task.Users, Data: task.Data, Constraints: task.Constraints,
		ExpectedResult: task.ExpectedResult, SuccessCriteria: task.SuccessCriteria,
		Contact: task.Contact, InteractionFormat: task.InteractionFormat,
	}
}

// RecalculateReadiness updates only derived fields in memory. The main API must
// verify human confirmation and save the card plus these fields atomically.
// It can also run on a copy of a draft for a preview; it does not confirm a task.
func (task *Task) RecalculateReadiness() {
	result := scoring.Calculate(task.Card())
	task.ReadinessScore = int16(result.Score)
	task.ReadinessLevel = strings.ToLower(result.Level)
	// Category contains only strings and integers, so JSON marshaling cannot fail.
	task.ReadinessBreakdown, _ = json.Marshal(result.Breakdown)
	task.MissingInformation = result.MissingFields
	task.Suggestions = make([]string, 0, len(result.Suggestions))
	for _, suggestion := range result.Suggestions {
		task.Suggestions = append(task.Suggestions, suggestion.Text)
	}
}
