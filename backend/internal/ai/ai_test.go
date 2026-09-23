package ai

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

type providerFunc func(context.Context, string, json.RawMessage, json.RawMessage) ([]byte, error)

func (f providerFunc) Complete(ctx context.Context, prompt string, input, schema json.RawMessage) ([]byte, error) {
	return f(ctx, prompt, input, schema)
}

func response(raw string) Provider {
	return providerFunc(func(_ context.Context, _ string, _, _ json.RawMessage) ([]byte, error) { return []byte(raw), nil })
}

const validQuestions = `{"questions":[{"field":"need","question":"Что изменить?"},{"field":"users","question":"Кто пользователи?"},{"field":"data","question":"Какие данные доступны?"}]}`

func TestClarifyValidation(t *testing.T) {
	for _, tc := range []struct {
		name, raw string
		valid     bool
	}{
		{"valid", validQuestions, true},
		{"broken JSON", `{`, false},
		{"trailing", validQuestions + `{}`, false},
		{"null", `null`, false},
		{"missing", `{}`, false},
		{"wrong type", `{"questions":42}`, false},
		{"too few", `{"questions":[{"field":"need","question":"Что изменить?"}]}`, false},
		{"null questions", `{"questions":null}`, false},
		{"question type", strings.Replace(validQuestions, `"Что изменить?"`, `3`, 1), false},
		{"empty question", strings.Replace(validQuestions, `Что изменить?`, `  `, 1), false},
		{"null question", strings.Replace(validQuestions, `"Что изменить?"`, `null`, 1), false},
		{"unknown field", strings.Replace(validQuestions, `"need"`, `"budget"`, 1), false},
		{"duplicate question", strings.Replace(validQuestions, `Кто пользователи?`, `Что изменить?`, 1), false},
		{"duplicate field", strings.Replace(validQuestions, `"users"`, `"need"`, 1), false},
		{"missing nested key", strings.Replace(validQuestions, `"field":"need",`, ``, 1), false},
		{"extra key", strings.Replace(validQuestions, `"field":"need"`, `"field":"need","extra":true`, 1), false},
		{"duplicate key", strings.Replace(validQuestions, `"field":"need"`, `"field":"need","field":"users"`, 1), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(response(tc.raw), 0, false).Clarify(context.Background(), "Нужно улучшить учёт")
			if tc.valid {
				if err != nil || got.Mode != "provider" || len(got.Questions) != 3 {
					t.Fatalf("%+v %v", got, err)
				}
			} else if !errors.Is(err, ErrResponse) {
				t.Fatalf("expected invalid response, got %v", err)
			}
		})
	}
}

func TestGenerateValidation(t *testing.T) {
	description := "Сейчас заказы в таблице."
	answer := "Сократить ручной ввод"
	input := GenerateInput{Description: description, Answers: []Answer{{Field: "need", Answer: answer}}}
	card := taskcard.Card{Context: &description, Need: &answer}
	valid, _ := json.Marshal(card)
	for _, tc := range []struct {
		name, raw string
		valid     bool
	}{
		{"valid", string(valid), true},
		{"missing fields", `{}`, false},
		{"broken", `{`, false},
		{"array", `[]`, false},
		{"wrong type", strings.Replace(string(valid), `"users":null`, `"users":[]`, 1), false},
		{"invented contact", strings.Replace(string(valid), `"contact":null`, `"contact":"invented@example.com"`, 1), false},
		{"answer wrong field", strings.Replace(string(valid), `"users":null`, `"users":"`+answer+`"`, 1), false},
		{"omitted answer", strings.Replace(string(valid), `"need":"`+answer+`"`, `"need":null`, 1), false},
		{"duplicate key", strings.Replace(string(valid), `"title":null`, `"title":null,"title":null`, 1), false},
		{"extra key", strings.Replace(string(valid), `"title":null`, `"title":null,"score":100`, 1), false},
		{"whitespace normalized", strings.Replace(string(valid), `"users":null`, `"users":"   "`, 1), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(response(tc.raw), 0, false).Generate(context.Background(), input)
			if tc.valid {
				if err != nil || got.Mode != "provider" || !reflect.DeepEqual(got.Card, card) {
					t.Fatalf("%+v %v", got, err)
				}
			} else if !errors.Is(err, ErrResponse) {
				t.Fatalf("expected invalid response, got %v", err)
			}
		})
	}
}

func TestFallbackAndFailures(t *testing.T) {
	timeout := providerFunc(func(ctx context.Context, _ string, _, _ json.RawMessage) ([]byte, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})
	failure := providerFunc(func(context.Context, string, json.RawMessage, json.RawMessage) ([]byte, error) {
		return nil, errors.New("secret provider body")
	})
	for _, tc := range []struct {
		name     string
		provider Provider
		want     error
		reason   string
	}{
		{"provider", failure, ErrProvider, "provider_error"},
		{"timeout", timeout, context.DeadlineExceeded, "timeout"},
		{"invalid", response(`{}`), ErrResponse, "invalid_response"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, fallback := range []bool{false, true} {
				s := New(tc.provider, time.Millisecond, fallback)
				q, err := s.Clarify(context.Background(), "Описание")
				g, genErr := s.Generate(context.Background(), GenerateInput{Description: "Описание"})
				if fallback {
					if err != nil || genErr != nil || q.Mode != "fallback" || g.Mode != "fallback" || q.FallbackReason != tc.reason || g.FallbackReason != tc.reason {
						t.Fatalf("%+v %+v %v %v", q, g, err, genErr)
					}
				} else if !errors.Is(err, tc.want) || !errors.Is(genErr, tc.want) {
					t.Fatalf("%v %v, expected %v", err, genErr, tc.want)
				}
			}
		})
	}
}

func TestInputAndCancellation(t *testing.T) {
	s := New(providerFunc(func(context.Context, string, json.RawMessage, json.RawMessage) ([]byte, error) {
		t.Fatal("provider should not be called")
		return nil, nil
	}), 0, true)
	for _, input := range []GenerateInput{
		{}, {Description: " \n"}, {Description: strings.Repeat("x", 20001)},
		{Description: "x", Answers: []Answer{{Field: "bad"}}},
		{Description: "x", Answers: []Answer{{Field: "need"}, {Field: "need"}}},
		{Description: "x", Answers: []Answer{{Field: "need", Answer: strings.Repeat("x", 20001)}}},
	} {
		if _, err := s.Generate(context.Background(), input); !errors.Is(err, ErrInput) {
			t.Fatalf("%v", err)
		}
	}
	if _, err := s.Clarify(context.Background(), " "); !errors.Is(err, ErrInput) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, service := range []*Service{s, New(nil, 0, false)} {
		if _, err := service.Clarify(ctx, "x"); !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
		if _, err := service.Generate(ctx, GenerateInput{Description: "x"}); !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	s = New(providerFunc(func(ctx context.Context, _ string, _, _ json.RawMessage) ([]byte, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}), time.Second, true)
	if _, err := s.Clarify(ctx, "x"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("parent timeout must not fall back: %v", err)
	}
}

func TestMockAndExample(t *testing.T) {
	raw, err := os.ReadFile("testdata/generate-input.json")
	if err != nil {
		t.Fatal(err)
	}
	var input GenerateInput
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}
	s := New(nil, 0, false)
	q, err := s.Clarify(context.Background(), input.Description)
	if err != nil || q.Mode != "mock" || len(q.Questions) < 3 {
		t.Fatalf("%+v %v", q, err)
	}
	qraw, _ := json.Marshal(struct {
		Questions []Question `json:"questions"`
	}{q.Questions})
	if _, err := validateQuestions(qraw); err != nil {
		t.Fatal(err)
	}
	got, err := s.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	wantRaw, err := os.ReadFile("testdata/generate-output.json")
	if err != nil {
		t.Fatal(err)
	}
	var want Generation
	if err := json.Unmarshal(wantRaw, &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
	cardRaw, _ := json.Marshal(got.Card)
	if _, err := validateCard(cardRaw, input); err != nil {
		t.Fatal(err)
	}
	again, _ := s.Generate(context.Background(), input)
	if !reflect.DeepEqual(got, again) {
		t.Fatal("mock not deterministic")
	}
	if score := scoring.Calculate(got.Card); score.Score != 60 {
		t.Fatalf("example score: %+v", score)
	}
	blank, err := s.Generate(context.Background(), GenerateInput{Description: "x", Answers: []Answer{{Field: "context", Answer: " "}}})
	if err != nil || blank.Card.Context != nil {
		t.Fatalf("blank answer: %+v %v", blank, err)
	}
}

func TestProviderReceivesPromptsAndSchema(t *testing.T) {
	s := New(providerFunc(func(_ context.Context, prompt string, input, schema json.RawMessage) ([]byte, error) {
		if prompt == "" || !json.Valid(input) || !json.Valid(schema) {
			t.Fatal("invalid request")
		}
		if prompt == clarifyPrompt {
			return []byte(validQuestions), nil
		}
		if prompt != generatePrompt {
			t.Fatal("unexpected prompt")
		}
		return json.Marshal(taskcard.Card{})
	}), 0, false)
	if _, err := s.Clarify(context.Background(), "x"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Generate(context.Background(), GenerateInput{Description: "x"}); err != nil {
		t.Fatal(err)
	}
}

func TestClarifyExample(t *testing.T) {
	raw, err := os.ReadFile("testdata/clarify-input.json")
	if err != nil {
		t.Fatal(err)
	}
	var input struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}
	raw, err = os.ReadFile("testdata/clarify-output.json")
	if err != nil {
		t.Fatal(err)
	}
	var want Clarification
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	modelResponse, _ := json.Marshal(struct {
		Questions []Question `json:"questions"`
	}{want.Questions})
	got, err := New(response(string(modelResponse)), 0, false).Clarify(context.Background(), input.Description)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("%+v %v", got, err)
	}
}

func TestExplicitAnswerOverridesDescription(t *testing.T) {
	description := "Срок 2 недели"
	for _, answer := range []string{"Срок 4 недели", " "} {
		input := GenerateInput{Description: description, Answers: []Answer{{Field: "constraints", Answer: answer}}}
		raw, _ := json.Marshal(taskcard.Card{Constraints: &description})
		if _, err := validateCard(raw, input); !errors.Is(err, ErrResponse) {
			t.Fatal("stale description must not override answer", err)
		}
	}
}
