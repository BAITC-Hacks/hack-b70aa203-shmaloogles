// Package scoring measures completeness, not semantic quality or factual accuracy.
package scoring

import "github.com/shmaloogles/business-task-platform/backend/internal/taskcard"

type Category struct {
	Name      string `json:"name"`
	Points    int    `json:"points"`
	MaxPoints int    `json:"max_points"`
}

type Suggestion struct {
	Field string `json:"field"`
	Text  string `json:"text"`
}

type Result struct {
	Score         int          `json:"score"`
	Level         string       `json:"level"`
	Breakdown     []Category   `json:"breakdown"`
	MissingFields []string     `json:"missing_fields"`
	Suggestions   []Suggestion `json:"suggestions"`
}

// Calculate is pure. The main API owns confirmation and saving/recalculating results.
func Calculate(card taskcard.Card) Result {
	r := Result{Breakdown: []Category{}, MissingFields: []string{}, Suggestions: []Suggestion{}}
	values := card.Values()
	groups := []struct {
		name   string
		weight int
		fields []string
	}{
		{"context_and_need", 20, []string{"context", "need"}},
		{"data", 20, []string{"data"}},
		{"expected_result", 15, []string{"expected_result"}},
		{"success_criteria", 15, []string{"success_criteria"}},
		{"constraints", 10, []string{"constraints"}},
		{"users", 10, []string{"users"}},
		{"business_communication", 10, []string{"contact", "interaction_format"}},
	}
	for _, group := range groups {
		category := Category{Name: group.name, MaxPoints: group.weight}
		for _, field := range group.fields {
			if taskcard.Filled(values[field]) {
				category.Points += group.weight / len(group.fields)
			}
		}
		r.Score += category.Points
		r.Breakdown = append(r.Breakdown, category)
	}
	hints := map[string]string{
		"title":              "Дайте задаче короткое название (не влияет на балл).",
		"topic":              "Укажите тему для фильтра каталога (не влияет на балл).",
		"context":            "Опишите текущий процесс и проблему: что происходит сейчас.",
		"need":               "Укажите, что бизнесу необходимо изменить в текущем процессе.",
		"users":              "Назовите пользователей решения и их рабочие задачи.",
		"data":               "Перечислите доступные данные, примеры или источники и условия доступа.",
		"constraints":        "Укажите сроки, допустимые технологии или ограничения доступа.",
		"expected_result":    "Укажите конкретный результат: например, прототип, отчёт или инструмент.",
		"success_criteria":   "Опишите измеримые условия приёмки результата.",
		"contact":            "Укажите контакт ответственного представителя бизнеса.",
		"interaction_format": "Опишите формат консультаций и порядок обратной связи.",
	}
	for _, field := range taskcard.Fields() {
		if !taskcard.Filled(values[field]) {
			r.MissingFields = append(r.MissingFields, field)
			r.Suggestions = append(r.Suggestions, Suggestion{Field: field, Text: hints[field]})
		}
	}
	r.Level = Level(r.Score)
	return r
}

// Level maps a readiness score in 0..100 to its product label.
func Level(score int) string {
	switch {
	case score < 40:
		return "Draft"
	case score < 70:
		return "Workable"
	case score < 90:
		return "Ready"
	default:
		return "Priority"
	}
}
