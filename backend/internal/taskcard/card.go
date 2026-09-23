// Package taskcard defines the shared AI/scoring input, without persistence metadata.
package taskcard

import "strings"

// Card uses nil for missing information. JSON always includes every field.
type Card struct {
	Title             *string `json:"title"`
	Topic             *string `json:"topic"`
	Context           *string `json:"context"`
	Need              *string `json:"need"`
	Users             *string `json:"users"`
	Data              *string `json:"data"`
	Constraints       *string `json:"constraints"`
	ExpectedResult    *string `json:"expected_result"`
	SuccessCriteria   *string `json:"success_criteria"`
	Contact           *string `json:"contact"`
	InteractionFormat *string `json:"interaction_format"`
}

// Fields returns stable schema order; callers receive their own slice.
func Fields() []string {
	return []string{"title", "topic", "context", "need", "users", "data", "constraints", "expected_result", "success_criteria", "contact", "interaction_format"}
}

func (c Card) Values() map[string]*string {
	return map[string]*string{
		"title": c.Title, "topic": c.Topic, "context": c.Context, "need": c.Need,
		"users": c.Users, "data": c.Data, "constraints": c.Constraints,
		"expected_result": c.ExpectedResult, "success_criteria": c.SuccessCriteria,
		"contact": c.Contact, "interaction_format": c.InteractionFormat,
	}
}

func Filled(value *string) bool { return value != nil && strings.TrimSpace(*value) != "" }
