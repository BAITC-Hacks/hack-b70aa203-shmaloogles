package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shmaloogles/business-task-platform/backend/internal/scoring"
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

func (store *Store) List(ctx context.Context, filter ListFilter) ([]Task, error) {
	query := `SELECT ` + taskColumns + ` FROM tasks WHERE status = 'published'`
	args := make([]any, 0, 2)
	if filter.Topic != "" {
		args = append(args, filter.Topic)
		query += fmt.Sprintf(" AND topic = $%d", len(args))
	}
	if filter.ReadinessLevel != "" {
		args = append(args, filter.ReadinessLevel)
		query += fmt.Sprintf(" AND readiness_level = $%d", len(args))
	}
	if filter.Sort == "readiness_asc" {
		query += " ORDER BY readiness_score ASC, published_at DESC"
	} else {
		query += " ORDER BY readiness_score DESC, published_at DESC"
	}

	rows, err := store.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list task rows: %w", err)
	}
	return tasks, nil
}

func (store *Store) Update(ctx context.Context, id int64, input UpdateInput, readiness scoring.Result) (Task, error) {
	breakdown, suggestions, err := readinessValues(readiness)
	if err != nil {
		return Task{}, err
	}

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
			readiness_score = $13,
			readiness_level = $14,
			readiness_breakdown = $15,
			missing_information = $16,
			suggestions = $17,
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
		readiness.Score,
		strings.ToLower(readiness.Level),
		breakdown,
		readiness.MissingFields,
		suggestions,
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

func (store *Store) Confirm(ctx context.Context, id int64, readiness scoring.Result) (Task, error) {
	breakdown, suggestions, err := readinessValues(readiness)
	if err != nil {
		return Task{}, err
	}

	row := store.db.QueryRow(ctx, `
		UPDATE tasks SET
			status = 'confirmed',
			readiness_score = $2,
			readiness_level = $3,
			readiness_breakdown = $4,
			missing_information = $5,
			suggestions = $6,
			confirmed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status = 'draft'
		RETURNING `+taskColumns,
		id,
		readiness.Score,
		strings.ToLower(readiness.Level),
		breakdown,
		readiness.MissingFields,
		suggestions,
	)
	return taskFromMutation(row, "confirm task")
}

func (store *Store) Publish(ctx context.Context, id int64) (Task, error) {
	row := store.db.QueryRow(ctx, `
		UPDATE tasks SET
			status = 'published',
			published_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND status = 'confirmed'
		RETURNING `+taskColumns,
		id,
	)
	return taskFromMutation(row, "publish task")
}

func taskFromMutation(row pgx.Row, operation string) (Task, error) {
	task, err := scanTask(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("%s: %w", operation, err)
	}
	return task, nil
}

func readinessValues(readiness scoring.Result) ([]byte, []string, error) {
	breakdown, err := json.Marshal(readiness.Breakdown)
	if err != nil {
		return nil, nil, fmt.Errorf("encode readiness breakdown: %w", err)
	}
	suggestions := make([]string, 0, len(readiness.Suggestions))
	for _, suggestion := range readiness.Suggestions {
		suggestions = append(suggestions, suggestion.Text)
	}
	return breakdown, suggestions, nil
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
