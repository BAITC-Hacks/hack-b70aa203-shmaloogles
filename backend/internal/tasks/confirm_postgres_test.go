package tasks

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
)

func TestConfirmUsesStoredCardInsteadOfStaleScore(t *testing.T) {
	dsn := os.Getenv("AI_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("requires migrated PostgreSQL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewStore(db)
	task, err := store.Create(ctx, CreateInput{InitialDescription: "confirmation regression test"})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", task.ID)
	stale := scoring.Calculate(task.Card())
	contextText := "Manual processing"
	card := UpdateInput{Context: &contextText}
	if _, err := store.Update(ctx, task.ID, card, scoring.Calculate(card)); err != nil {
		t.Fatal(err)
	}
	confirmed, err := store.Confirm(ctx, task.ID, stale)
	if err != nil {
		t.Fatal(err)
	}
	if confirmed.Status != "confirmed" || confirmed.ReadinessScore != 10 {
		t.Fatalf("confirmation saved stale score: status=%s score=%d", confirmed.Status, confirmed.ReadinessScore)
	}
}
