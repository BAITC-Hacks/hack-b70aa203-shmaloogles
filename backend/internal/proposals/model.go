package proposals

import "time"

type Proposal struct {
	ID           int64     `json:"id"`
	TaskID       int64     `json:"task_id"`
	TeamID       int64     `json:"team_id"`
	TeamName     string    `json:"team_name"`
	SolutionIdea string    `json:"solution_idea"`
	Plan         string    `json:"plan"`
	Timeline     string    `json:"timeline"`
	PrototypeURL *string   `json:"prototype_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateInput struct {
	TeamID       int64   `json:"team_id"`
	SolutionIdea string  `json:"solution_idea"`
	Plan         string  `json:"plan"`
	Timeline     string  `json:"timeline"`
	PrototypeURL *string `json:"prototype_url"`
}

type StatusInput struct {
	Status string `json:"status"`
}
