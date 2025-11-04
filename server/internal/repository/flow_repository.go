package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/visual-api-testing-platform/server/internal/models"
)

// FlowRepository handles flow database operations
type FlowRepository struct {
	db         *pgxpool.Pool
	nodeRepo   *NodeRepository
}

// NewFlowRepository creates a new flow repository
func NewFlowRepository(db *pgxpool.Pool) *FlowRepository {
	return &FlowRepository{
		db:       db,
		nodeRepo: NewNodeRepository(db),
	}
}

// Create creates a new flow
func (r *FlowRepository) Create(ctx context.Context, flow *models.Flow) error {
	edgesJSON, _ := json.Marshal(flow.Edges)

	query := `
		INSERT INTO flows (id, user_id, name, description, tags, edges, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		flow.ID,
		flow.UserID,
		flow.Name,
		flow.Description,
		flow.Tags,
		edgesJSON,
		flow.CreatedAt,
		flow.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// Create nodes separately
	for i := range flow.Nodes {
		if err := r.nodeRepo.Create(ctx, flow.ID, &flow.Nodes[i]); err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves a flow by ID with nodes joined
func (r *FlowRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Flow, error) {
	var flow models.Flow
	var edgesJSON []byte

	query := `
		SELECT id, user_id, name, description, tags, edges, created_at, updated_at
		FROM flows
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&flow.ID,
		&flow.UserID,
		&flow.Name,
		&flow.Description,
		&flow.Tags,
		&edgesJSON,
		&flow.CreatedAt,
		&flow.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse edges
	json.Unmarshal(edgesJSON, &flow.Edges)

	// Get nodes using JOIN
	nodes, err := r.nodeRepo.GetByFlowID(ctx, id)
	if err != nil {
		return nil, err
	}
	flow.Nodes = nodes

	return &flow, nil
}

// GetByUserID retrieves all flows for a user with nodes joined
func (r *FlowRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Flow, error) {
	query := `
		SELECT id, user_id, name, description, tags, edges, created_at, updated_at
		FROM flows
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flows []models.Flow
	for rows.Next() {
		var flow models.Flow
		var edgesJSON []byte

		err := rows.Scan(
			&flow.ID,
			&flow.UserID,
			&flow.Name,
			&flow.Description,
			&flow.Tags,
			&edgesJSON,
			&flow.CreatedAt,
			&flow.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse edges
		json.Unmarshal(edgesJSON, &flow.Edges)

		// Get nodes using JOIN
		nodes, err := r.nodeRepo.GetByFlowID(ctx, flow.ID)
		if err != nil {
			return nil, err
		}
		flow.Nodes = nodes

		flows = append(flows, flow)
	}

	return flows, nil
}

// Update updates a flow and its nodes
func (r *FlowRepository) Update(ctx context.Context, flow *models.Flow) error {
	edgesJSON, _ := json.Marshal(flow.Edges)

	query := `
		UPDATE flows
		SET name = $2, description = $3, tags = $4, edges = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		flow.ID,
		flow.Name,
		flow.Description,
		flow.Tags,
		edgesJSON,
		time.Now(),
	)

	if err != nil {
		return err
	}

	// Update nodes using bulk upsert
	return r.nodeRepo.BulkUpsert(ctx, flow.ID, flow.Nodes)
}

// Delete deletes a flow (nodes will be deleted via CASCADE)
func (r *FlowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM flows WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

