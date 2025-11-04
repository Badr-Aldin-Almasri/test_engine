package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/visual-api-testing-platform/server/internal/models"
)

// TestRunRepository handles test run database operations
type TestRunRepository struct {
	db *pgxpool.Pool
}

// NewTestRunRepository creates a new test run repository
func NewTestRunRepository(db *pgxpool.Pool) *TestRunRepository {
	return &TestRunRepository{db: db}
}

// Create creates a new test run
func (r *TestRunRepository) Create(ctx context.Context, testRun *models.TestRun) error {
	nodeResultsJSON, _ := json.Marshal(testRun.NodeResults)

	query := `
		INSERT INTO test_runs (id, flow_id, status, started_at, completed_at, duration_ms, node_results, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		testRun.ID,
		testRun.FlowID,
		string(testRun.Status),
		testRun.StartedAt,
		testRun.CompletedAt,
		testRun.DurationMs,
		nodeResultsJSON,
		testRun.Error,
	)

	return err
}

// GetByID retrieves a test run by ID
func (r *TestRunRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TestRun, error) {
	var testRun models.TestRun
	var statusStr string
	var nodeResultsJSON []byte

	query := `
		SELECT id, flow_id, status, started_at, completed_at, duration_ms, node_results, error, created_at
		FROM test_runs
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&testRun.ID,
		&testRun.FlowID,
		&statusStr,
		&testRun.StartedAt,
		&testRun.CompletedAt,
		&testRun.DurationMs,
		&nodeResultsJSON,
		&testRun.Error,
		&testRun.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	testRun.Status = models.ExecutionStatus(statusStr)
	json.Unmarshal(nodeResultsJSON, &testRun.NodeResults)

	return &testRun, nil
}

// GetByFlowID retrieves all test runs for a flow
func (r *TestRunRepository) GetByFlowID(ctx context.Context, flowID uuid.UUID, limit int) ([]models.TestRun, error) {
	query := `
		SELECT id, flow_id, status, started_at, completed_at, duration_ms, node_results, error, created_at
		FROM test_runs
		WHERE flow_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, flowID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testRuns []models.TestRun
	for rows.Next() {
		var testRun models.TestRun
		var statusStr string
		var nodeResultsJSON []byte

		err := rows.Scan(
			&testRun.ID,
			&testRun.FlowID,
			&statusStr,
			&testRun.StartedAt,
			&testRun.CompletedAt,
			&testRun.DurationMs,
			&nodeResultsJSON,
			&testRun.Error,
			&testRun.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		testRun.Status = models.ExecutionStatus(statusStr)
		json.Unmarshal(nodeResultsJSON, &testRun.NodeResults)

		testRuns = append(testRuns, testRun)
	}

	return testRuns, nil
}

// Update updates a test run
func (r *TestRunRepository) Update(ctx context.Context, testRun *models.TestRun) error {
	nodeResultsJSON, _ := json.Marshal(testRun.NodeResults)

	query := `
		UPDATE test_runs
		SET status = $2, completed_at = $3, duration_ms = $4, node_results = $5, error = $6
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		testRun.ID,
		string(testRun.Status),
		testRun.CompletedAt,
		testRun.DurationMs,
		nodeResultsJSON,
		testRun.Error,
	)

	return err
}

