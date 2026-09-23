package proposals

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound        = errors.New("proposal not found")
	ErrTaskNotFound    = errors.New("task not found")
	ErrTaskUnavailable = errors.New("task is not published")
	ErrTeamNotFound    = errors.New("team not found")
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (store *Store) Create(ctx context.Context, taskID int64, input CreateInput) (Proposal, error) {
	var taskExists, taskPublished, teamExists bool
	if err := store.db.QueryRow(ctx, `
		SELECT
			EXISTS (SELECT 1 FROM tasks WHERE id = $1),
			EXISTS (SELECT 1 FROM tasks WHERE id = $1 AND status = 'published'),
			EXISTS (SELECT 1 FROM teams WHERE id = $2)`,
		taskID, input.TeamID,
	).Scan(&taskExists, &taskPublished, &teamExists); err != nil {
		return Proposal{}, fmt.Errorf("validate proposal references: %w", err)
	}
	if !taskExists {
		return Proposal{}, ErrTaskNotFound
	}
	if !taskPublished {
		return Proposal{}, ErrTaskUnavailable
	}
	if !teamExists {
		return Proposal{}, ErrTeamNotFound
	}

	row := store.db.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO proposals (task_id, team_id, solution_idea, plan, timeline, prototype_url)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING *
		)
		SELECT `+proposalColumns+`
		FROM inserted p
		JOIN teams t ON t.id = p.team_id`,
		taskID, input.TeamID, input.SolutionIdea, input.Plan, input.Timeline, input.PrototypeURL,
	)
	proposal, err := scanProposal(row)
	if err != nil {
		return Proposal{}, fmt.Errorf("create proposal: %w", err)
	}
	return proposal, nil
}

func (store *Store) ListByTask(ctx context.Context, taskID int64) ([]Proposal, error) {
	var taskExists bool
	if err := store.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)`, taskID).Scan(&taskExists); err != nil {
		return nil, fmt.Errorf("validate task: %w", err)
	}
	if !taskExists {
		return nil, ErrTaskNotFound
	}

	rows, err := store.db.Query(ctx, `
		SELECT `+proposalColumns+`
		FROM proposals p
		JOIN teams t ON t.id = p.team_id
		WHERE p.task_id = $1
		ORDER BY p.created_at DESC`,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("list proposals: %w", err)
	}
	defer rows.Close()

	items := make([]Proposal, 0)
	for rows.Next() {
		proposal, err := scanProposal(rows)
		if err != nil {
			return nil, fmt.Errorf("scan proposal: %w", err)
		}
		items = append(items, proposal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list proposal rows: %w", err)
	}
	return items, nil
}

func (store *Store) UpdateStatus(ctx context.Context, id int64, status string) (Proposal, error) {
	row := store.db.QueryRow(ctx, `
		WITH updated AS (
			UPDATE proposals SET status = $2, updated_at = NOW()
			WHERE id = $1
			RETURNING *
		)
		SELECT `+proposalColumns+`
		FROM updated p
		JOIN teams t ON t.id = p.team_id`,
		id, status,
	)
	proposal, err := scanProposal(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Proposal{}, ErrNotFound
	}
	if err != nil {
		return Proposal{}, fmt.Errorf("update proposal status: %w", err)
	}
	return proposal, nil
}

const proposalColumns = `
	p.id, p.task_id, p.team_id, t.name, p.solution_idea, p.plan, p.timeline,
	p.prototype_url, p.status, p.created_at, p.updated_at`

func scanProposal(row pgx.Row) (Proposal, error) {
	var proposal Proposal
	err := row.Scan(
		&proposal.ID,
		&proposal.TaskID,
		&proposal.TeamID,
		&proposal.TeamName,
		&proposal.SolutionIdea,
		&proposal.Plan,
		&proposal.Timeline,
		&proposal.PrototypeURL,
		&proposal.Status,
		&proposal.CreatedAt,
		&proposal.UpdatedAt,
	)
	return proposal, err
}
