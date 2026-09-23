package tasks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// RecalculateAll updates only derived readiness fields. Without apply it rolls
// back the transaction. Row locks keep concurrent card edits consistent.
func (store *Store) RecalculateAll(ctx context.Context, apply bool) (int, error) {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())
	count, err := recalculateAll(ctx, tx)
	if err != nil {
		return 0, err
	}
	if apply {
		if err := tx.Commit(ctx); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func recalculateAll(ctx context.Context, tx pgx.Tx) (int, error) {
	rows, err := tx.Query(ctx, `SELECT `+taskColumns+` FROM tasks ORDER BY id FOR UPDATE`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return 0, err
		}
		items = append(items, task)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	rows.Close()
	for _, task := range items {
		task.RecalculateReadiness()
		_, err := tx.Exec(ctx, `UPDATE tasks SET readiness_score=$2, readiness_level=$3,
			readiness_breakdown=$4, missing_information=$5, suggestions=$6 WHERE id=$1`,
			task.ID, task.ReadinessScore, task.ReadinessLevel, task.ReadinessBreakdown,
			task.MissingInformation, task.Suggestions)
		if err != nil {
			return 0, fmt.Errorf("recalculate task %d: %w", task.ID, err)
		}
	}
	return len(items), nil
}
