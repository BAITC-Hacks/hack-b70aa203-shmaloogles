package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (store *Store) Create(ctx context.Context, input CreateInput) (Task, error) {
	row := store.db.QueryRow(ctx, `
		INSERT INTO tasks (initial_description)
		VALUES ($1)
		RETURNING `+taskColumns,
		input.InitialDescription,
	)

	task, err := scanTask(row)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (store *Store) Get(ctx context.Context, id int64) (Task, error) {
	row := store.db.QueryRow(ctx, `SELECT `+taskColumns+` FROM tasks WHERE id = $1`, id)
	task, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("get task: %w", err)
	}
	return task, nil
}

func (store *Store) Update(ctx context.Context, id int64, input UpdateInput) (Task, error) {
	row := store.db.QueryRow(ctx, `
		UPDATE tasks SET
			title = $2,
			topic = $3,
			context = $4,
			need = $5,
			users = $6,
			data = $7,
			constraints = $8,
			expected_result = $9,
			success_criteria = $10,
			contact = $11,
			interaction_format = $12,
			updated_at = NOW()
		WHERE id = $1
		RETURNING `+taskColumns,
		id,
		input.Title,
		input.Topic,
		input.Context,
		input.Need,
		input.Users,
		input.Data,
		input.Constraints,
		input.ExpectedResult,
		input.SuccessCriteria,
		input.Contact,
		input.InteractionFormat,
	)

	task, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	return task, nil
}

const taskColumns = `
	id, initial_description, clarification, title, topic, context, need, users,
	data, constraints, expected_result, success_criteria, contact, interaction_format,
	status, readiness_score, readiness_level, readiness_breakdown,
	missing_information, suggestions, confirmed_at, published_at, created_at, updated_at`

func scanTask(row pgx.Row) (Task, error) {
	var task Task
	err := row.Scan(
		&task.ID,
		&task.InitialDescription,
		&task.Clarification,
		&task.Title,
		&task.Topic,
		&task.Context,
		&task.Need,
		&task.Users,
		&task.Data,
		&task.Constraints,
		&task.ExpectedResult,
		&task.SuccessCriteria,
		&task.Contact,
		&task.InteractionFormat,
		&task.Status,
		&task.ReadinessScore,
		&task.ReadinessLevel,
		&task.ReadinessBreakdown,
		&task.MissingInformation,
		&task.Suggestions,
		&task.ConfirmedAt,
		&task.PublishedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}
