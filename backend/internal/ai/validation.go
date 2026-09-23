package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

// object rejects duplicate/unknown/missing keys, trailing JSON and non-objects.
func object(raw []byte, fields []string) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(raw))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return nil, ErrResponse
	}
	result := map[string]json.RawMessage{}
	allowed := map[string]bool{}
	for _, field := range fields {
		allowed[field] = true
	}
	for d.More() {
		token, err := d.Token()
		if err != nil {
			return nil, ErrResponse
		}
		key, ok := token.(string)
		if !ok || !allowed[key] || result[key] != nil {
			return nil, ErrResponse
		}
		var value json.RawMessage
		if d.Decode(&value) != nil {
			return nil, ErrResponse
		}
		result[key] = value
	}
	if _, err := d.Token(); err != nil {
		return nil, ErrResponse
	}
	if _, err := d.Token(); err != io.EOF {
		return nil, ErrResponse
	}
	if len(result) != len(fields) {
		return nil, ErrResponse
	}
	return result, nil
}

func validateQuestions(raw []byte) ([]Question, error) {
	obj, err := object(raw, []string{"questions"})
	if err != nil {
		return nil, err
	}
	var rows []json.RawMessage
	if json.Unmarshal(obj["questions"], &rows) != nil || len(rows) < 3 || len(rows) > 11 {
		return nil, ErrResponse
	}
	questions := make([]Question, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		if _, err := object(row, []string{"field", "question"}); err != nil {
			return nil, err
		}
		var q Question
		if json.Unmarshal(row, &q) != nil {
			return nil, ErrResponse
		}
		q.Question = strings.TrimSpace(q.Question)
		key := strings.ToLower(q.Question)
		if !validField(q.Field) || q.Question == "" || len(q.Question) > 2000 || seen[key] || seen["field:"+q.Field] {
			return nil, ErrResponse
		}
		seen[key], seen["field:"+q.Field] = true, true
		questions = append(questions, q)
	}
	return questions, nil
}

func validateCard(raw []byte, input GenerateInput) (taskcard.Card, error) {
	obj, err := object(raw, taskcard.Fields())
	if err != nil {
		return taskcard.Card{}, err
	}
	answers := map[string]string{}
	for _, a := range input.Answers {
		answers[a.Field] = strings.TrimSpace(a.Answer)
	}
	for _, field := range taskcard.Fields() {
		var value *string
		if json.Unmarshal(obj[field], &value) != nil {
			return taskcard.Card{}, ErrResponse
		}
		answer, answered := answers[field]
		if value != nil {
			trimmed := strings.TrimSpace(*value)
			if trimmed == "" {
				value = nil
			} else {
				if (answered && trimmed != answer) || (!answered && !strings.Contains(input.Description, trimmed)) {
					return taskcard.Card{}, fmt.Errorf("%w: ungrounded field %s", ErrResponse, field)
				}
				value = &trimmed
			}
		}
		if answered && answer != "" && value == nil {
			return taskcard.Card{}, fmt.Errorf("%w: omitted answer %s", ErrResponse, field)
		}
		obj[field], _ = json.Marshal(value)
	}
	normalized, _ := json.Marshal(obj)
	var card taskcard.Card
	if json.Unmarshal(normalized, &card) != nil {
		return taskcard.Card{}, ErrResponse
	}
	return card, nil
}

func cardSchema() json.RawMessage {
	properties := map[string]any{}
	for _, field := range taskcard.Fields() {
		properties[field] = map[string]any{"type": []string{"string", "null"}}
	}
	schema, _ := json.Marshal(map[string]any{"type": "object", "properties": properties, "required": taskcard.Fields(), "additionalProperties": false})
	return schema
}

func questionSchema() json.RawMessage {
	schema, _ := json.Marshal(map[string]any{
		"type": "object", "required": []string{"questions"}, "additionalProperties": false,
		"properties": map[string]any{"questions": map[string]any{
			"type": "array", "minItems": 3, "maxItems": 11,
			"items": map[string]any{"type": "object", "required": []string{"field", "question"}, "additionalProperties": false,
				"properties": map[string]any{"field": map[string]any{"type": "string", "enum": taskcard.Fields()}, "question": map[string]any{"type": "string"}}},
		}},
	})
	return schema
}
