package scoring

import (
	"reflect"
	"testing"

	"github.com/shmaloogles/business-task-platform/backend/internal/taskcard"
)

func str(s string) *string { return &s }

func TestCalculate(t *testing.T) {
	full := taskcard.Card{Title: str("T"), Topic: str("T"), Context: str("C"), Need: str("N"), Users: str("U"), Data: str("D"), Constraints: str("C"), ExpectedResult: str("E"), SuccessCriteria: str("S"), Contact: str("C"), InteractionFormat: str("I")}
	cases := []struct {
		name    string
		card    taskcard.Card
		score   int
		level   string
		missing int
	}{
		{"empty", taskcard.Card{}, 0, "Draft", 11},
		{"full", full, 100, "Priority", 0},
		{"context only", taskcard.Card{Context: str("context")}, 10, "Draft", 10},
		{"need only", taskcard.Card{Need: str("need")}, 10, "Draft", 10},
		{"context and need", taskcard.Card{Context: str("context"), Need: str("need")}, 20, "Draft", 9},
		{"data", taskcard.Card{Data: str("data")}, 20, "Draft", 10},
		{"result", taskcard.Card{ExpectedResult: str("result")}, 15, "Draft", 10},
		{"criteria", taskcard.Card{SuccessCriteria: str("criteria")}, 15, "Draft", 10},
		{"constraints", taskcard.Card{Constraints: str("constraints")}, 10, "Draft", 10},
		{"users", taskcard.Card{Users: str("users")}, 10, "Draft", 10},
		{"contact only", taskcard.Card{Contact: str("contact")}, 5, "Draft", 10},
		{"format only", taskcard.Card{InteractionFormat: str("format")}, 5, "Draft", 10},
		{"communication", taskcard.Card{Contact: str("contact"), InteractionFormat: str("format")}, 10, "Draft", 9},
		{"unscored fields", taskcard.Card{Title: str("title"), Topic: str("topic")}, 0, "Draft", 9},
		{"whitespace", taskcard.Card{Context: str(" \t\n"), Need: str("\u00a0"), Data: str("")}, 0, "Draft", 11},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := Calculate(tc.card)
			if r.Score != tc.score || r.Level != tc.level || len(r.MissingFields) != tc.missing {
				t.Fatalf("unexpected result: %+v", r)
			}
			if len(r.Breakdown) != 7 || len(r.Suggestions) != tc.missing {
				t.Fatalf("incomplete breakdown: %+v", r)
			}
			total, max := 0, 0
			for _, c := range r.Breakdown {
				total += c.Points
				max += c.MaxPoints
				if c.Points < 0 || c.Points > c.MaxPoints {
					t.Fatal(c)
				}
			}
			if total != r.Score || max != 100 {
				t.Fatalf("bad totals: %+v", r)
			}
			for i, hint := range r.Suggestions {
				if hint.Field != r.MissingFields[i] || hint.Text == "" {
					t.Fatal(hint)
				}
			}
			if !reflect.DeepEqual(r, Calculate(tc.card)) {
				t.Fatal("not deterministic")
			}
		})
	}
}

func TestLevelBoundaries(t *testing.T) {
	for score, want := range map[int]string{0: "Draft", 39: "Draft", 40: "Workable", 69: "Workable", 70: "Ready", 89: "Ready", 90: "Priority", 100: "Priority"} {
		if got := Level(score); got != want {
			t.Errorf("Level(%d)=%s, want %s", score, got, want)
		}
	}
}

func TestRecalculateAfterEdit(t *testing.T) {
	card := taskcard.Card{Context: str("process")}
	before := Calculate(card)
	card.Need = str("change")
	card.Data = str("CSV")
	after := Calculate(card)
	if before.Score != 10 || after.Score != 40 || after.Level != "Workable" {
		t.Fatalf("%+v -> %+v", before, after)
	}
	card.Data = str(" ")
	if Calculate(card).Score != 20 {
		t.Fatal("removed data still scores")
	}
}
