package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRecalculateAllPostgres(t *testing.T) {
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
	// A private schema prevents any modification of existing demo/user tasks.
	schema := fmt.Sprintf("recalculate_test_%d", time.Now().UnixNano())
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := db.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := db.Exec(context.Background(), "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
	}()
	if _, err := db.Exec(ctx, "CREATE TABLE "+quoted+".tasks (LIKE public.tasks INCLUDING ALL)"); err != nil {
		t.Fatal(err)
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	testDB, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()
	store := NewStore(testDB)
	for _, status := range []string{"draft", "confirmed", "published"} {
		var id int64
		if err := testDB.QueryRow(ctx, `INSERT INTO tasks (initial_description, context, status, readiness_score, readiness_level, confirmed_at, published_at)
			VALUES ('test', 'manual process', $1, 99, 'priority', CASE WHEN $1 <> 'draft' THEN NOW() END, CASE WHEN $1 = 'published' THEN NOW() END) RETURNING id`, status).Scan(&id); err != nil {
			t.Fatal(err)
		}
		before, err := store.Get(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.RecalculateAll(ctx, false); err != nil {
			t.Fatal(err)
		}
		after, err := store.Get(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(before, after) {
			t.Fatal("dry-run changed task")
		}
		if _, err := store.RecalculateAll(ctx, true); err != nil {
			t.Fatal(err)
		}
		after, err = store.Get(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		before.RecalculateReadiness()
		var expectedBreakdown, actualBreakdown any
		if err := json.Unmarshal(before.ReadinessBreakdown, &expectedBreakdown); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(after.ReadinessBreakdown, &actualBreakdown); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expectedBreakdown, actualBreakdown) {
			t.Fatal("incorrect breakdown")
		}
		// JSONB formatting differs from json.Marshal; compare its decoded value separately.
		before.ReadinessBreakdown = after.ReadinessBreakdown
		if !reflect.DeepEqual(before, after) {
			t.Fatal("apply changed non-derived fields or produced incorrect readiness")
		}
		if _, err := store.RecalculateAll(ctx, true); err != nil {
			t.Fatal(err)
		}
		again, err := store.Get(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(after, again) {
			t.Fatal("recalculation is not idempotent")
		}
	}
}
