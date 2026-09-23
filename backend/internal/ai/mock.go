package ai

import (
	"encoding/json"
	"strings"

	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

// The mock deliberately does not attempt semantic extraction from free text.
func mockClarify() Clarification {
	return Clarification{Metadata: Metadata{Mode: "mock"}, Questions: []Question{
		{Field: "need", Question: "Что именно необходимо изменить в описанной ситуации?"},
		{Field: "users", Question: "Кто будет пользоваться решением и для каких действий?"},
		{Field: "data", Question: "Какие данные, примеры или материалы доступны команде и как получить доступ?"},
		{Field: "constraints", Question: "Какие сроки, технологии и ограничения доступа нужно учитывать?"},
		{Field: "expected_result", Question: "Какой конкретный результат должна передать команда?"},
		{Field: "success_criteria", Question: "По каким измеримым признакам вы примете результат?"},
		{Field: "contact", Question: "Кто представляет бизнес и как с ним связаться?"},
		{Field: "interaction_format", Question: "Как будут проходить консультации и обратная связь?"},
		{Field: "title", Question: "Как вы хотите назвать эту задачу?"},
		{Field: "topic", Question: "К какой теме относится задача для фильтра каталога?"},
	}}
}

func mockGenerate(input GenerateInput) Generation {
	values := map[string]*string{}
	description := strings.TrimSpace(input.Description)
	values["context"] = &description
	for _, answer := range input.Answers {
		value := strings.TrimSpace(answer.Answer)
		if value == "" {
			values[answer.Field] = nil
		} else {
			values[answer.Field] = &value
		}
	}
	raw, _ := json.Marshal(values)
	var card taskcard.Card
	_ = json.Unmarshal(raw, &card)
	return Generation{Card: card, Metadata: Metadata{Mode: "mock"}}
}
