// Package ai clarifies descriptions and generates validated, extractive task cards.
package ai

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

//go:embed prompts/clarify.txt
var clarifyPrompt string

//go:embed prompts/generate.txt
var generatePrompt string

var (
	ErrInput    = errors.New("invalid AI input")
	ErrResponse = errors.New("invalid AI response")
	ErrProvider = errors.New("AI provider failed")
)

type Question struct {
	Field    string `json:"field"`
	Question string `json:"question"`
}
type Answer struct {
	Field  string `json:"field"`
	Answer string `json:"answer"`
}
type GenerateInput struct {
	Description string   `json:"description"`
	Answers     []Answer `json:"answers"`
}
type Metadata struct {
	Mode           string `json:"mode"` // provider, mock or fallback
	FallbackReason string `json:"fallback_reason,omitempty"`
}
type Clarification struct {
	Questions []Question `json:"questions"`
	Metadata
}
type Generation struct {
	Card taskcard.Card `json:"card"`
	Metadata
}

// Provider implementations must respect ctx. Output is raw model JSON, not an envelope.
type Provider interface {
	Complete(ctx context.Context, prompt string, input json.RawMessage, schema json.RawMessage) ([]byte, error)
}

type Service struct {
	provider Provider
	timeout  time.Duration
	fallback bool
}

// New uses an explicit mock when provider is nil. Fallback is opt-in for provider errors.
func New(provider Provider, timeout time.Duration, fallback bool) *Service {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &Service{provider: provider, timeout: timeout, fallback: fallback}
}

func (s *Service) Clarify(ctx context.Context, description string) (Clarification, error) {
	if err := validateInput(GenerateInput{Description: description}); err != nil {
		return Clarification{}, err
	}
	if err := ctx.Err(); err != nil {
		return Clarification{}, err
	}
	if s.provider == nil {
		return mockClarify(), nil
	}
	input, _ := json.Marshal(struct {
		Description string `json:"description"`
	}{description})
	raw, err := s.complete(ctx, clarifyPrompt, input, questionSchema())
	var questions []Question
	if err == nil {
		questions, err = validateQuestions(raw)
	}
	if err != nil {
		if ctx.Err() != nil {
			return Clarification{}, ctx.Err()
		}
		if !s.fallback {
			return Clarification{}, err
		}
		result := mockClarify()
		result.Metadata = Metadata{Mode: "fallback", FallbackReason: reason(err)}
		return result, nil
	}
	return Clarification{Questions: questions, Metadata: Metadata{Mode: "provider"}}, nil
}

func (s *Service) Generate(ctx context.Context, input GenerateInput) (Generation, error) {
	if err := validateInput(input); err != nil {
		return Generation{}, err
	}
	if err := ctx.Err(); err != nil {
		return Generation{}, err
	}
	if s.provider == nil {
		return mockGenerate(input), nil
	}
	encoded, _ := json.Marshal(input)
	raw, err := s.complete(ctx, generatePrompt, encoded, cardSchema())
	var card taskcard.Card
	if err == nil {
		card, err = validateCard(raw, input)
	}
	if err != nil {
		if ctx.Err() != nil {
			return Generation{}, ctx.Err()
		}
		if !s.fallback {
			return Generation{}, err
		}
		result := mockGenerate(input)
		result.Metadata = Metadata{Mode: "fallback", FallbackReason: reason(err)}
		return result, nil
	}
	return Generation{Card: card, Metadata: Metadata{Mode: "provider"}}, nil
}

func (s *Service) complete(ctx context.Context, prompt string, input, schema json.RawMessage) ([]byte, error) {
	child, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	raw, err := s.provider.Complete(child, prompt, input, schema)
	if child.Err() != nil {
		return nil, child.Err()
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, context.DeadlineExceeded
		}
		if errors.Is(err, ErrResponse) {
			return nil, ErrResponse
		}
		// Never propagate provider bodies, credentials or request URLs to API clients.
		return nil, ErrProvider
	}
	return raw, nil
}

func reason(err error) string {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	case errors.Is(err, ErrResponse):
		return "invalid_response"
	default:
		return "provider_error"
	}
}

func validateInput(input GenerateInput) error {
	if strings.TrimSpace(input.Description) == "" || len(input.Description) > 20000 {
		return fmt.Errorf("%w: description must contain 1..20000 bytes", ErrInput)
	}
	seen := map[string]bool{}
	for _, answer := range input.Answers {
		if !validField(answer.Field) || seen[answer.Field] || len(answer.Answer) > 20000 {
			return fmt.Errorf("%w: answers require unique known fields and at most 20000 bytes each", ErrInput)
		}
		seen[answer.Field] = true
	}
	return nil
}

func validField(field string) bool {
	for _, known := range taskcard.Fields() {
		if known == field {
			return true
		}
	}
	return false
}
