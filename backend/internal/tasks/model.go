package tasks

import (
	"encoding/json"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

type Task struct {
	ID                 int64           `json:"id"`
	InitialDescription *string         `json:"initial_description"`
	Clarification      json.RawMessage `json:"clarification"`
	Title              *string         `json:"title"`
	Topic              *string         `json:"topic"`
	Context            *string         `json:"context"`
	Need               *string         `json:"need"`
	Users              *string         `json:"users"`
	Data               *string         `json:"data"`
	Constraints        *string         `json:"constraints"`
	ExpectedResult     *string         `json:"expected_result"`
	SuccessCriteria    *string         `json:"success_criteria"`
	Contact            *string         `json:"contact"`
	InteractionFormat  *string         `json:"interaction_format"`
	Status             string          `json:"status"`
	ReadinessScore     int16           `json:"readiness_score"`
	ReadinessLevel     string          `json:"readiness_level"`
	ReadinessBreakdown json.RawMessage `json:"readiness_breakdown"`
	MissingInformation []string        `json:"missing_information"`
	Suggestions        []string        `json:"suggestions"`
	ConfirmedAt        *time.Time      `json:"confirmed_at"`
	PublishedAt        *time.Time      `json:"published_at"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type CreateInput struct {
	InitialDescription string `json:"initial_description"`
}

type UpdateInput = taskcard.Card

func (task Task) Card() taskcard.Card {
	return taskcard.Card{
		Title:             task.Title,
		Topic:             task.Topic,
		Context:           task.Context,
		Need:              task.Need,
		Users:             task.Users,
		Data:              task.Data,
		Constraints:       task.Constraints,
		ExpectedResult:    task.ExpectedResult,
		SuccessCriteria:   task.SuccessCriteria,
		Contact:           task.Contact,
		InteractionFormat: task.InteractionFormat,
	}
}
