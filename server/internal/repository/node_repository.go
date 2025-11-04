package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/visual-api-testing-platform/server/internal/models"
)

// NodeRepository handles node database operations
type NodeRepository struct {
	db *pgxpool.Pool
}

// NewNodeRepository creates a new node repository
func NewNodeRepository(db *pgxpool.Pool) *NodeRepository {
	return &NodeRepository{db: db}
}

// Create creates a new node
func (r *NodeRepository) Create(ctx context.Context, flowID uuid.UUID, node *models.FlowNode) error {
	configJSON, _ := json.Marshal(node.Data.Config)
	outputJSON, _ := json.Marshal(node.Data.Output)

	query := `
		INSERT INTO flow_nodes (
			flow_id, node_id, type, position_x, position_y,
			data_id, data_type, data_label, data_status,
			data_config, data_output, data_error,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		flowID,
		node.ID,
		node.Type,
		node.Position.X,
		node.Position.Y,
		node.Data.ID,
		node.Data.Type,
		node.Data.Label,
		node.Data.Status,
		configJSON,
		outputJSON,
		node.Data.Error,
		time.Now(),
		time.Now(),
	)

	return err
}

// GetByFlowID retrieves all nodes for a flow
func (r *NodeRepository) GetByFlowID(ctx context.Context, flowID uuid.UUID) ([]models.FlowNode, error) {
	query := `
		SELECT 
			node_id, type, position_x, position_y,
			data_id, data_type, data_label, data_status,
			data_config, data_output, data_error
		FROM flow_nodes
		WHERE flow_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, flowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []models.FlowNode
	for rows.Next() {
		var node models.FlowNode
		var configJSON, outputJSON []byte
		var status *string
		var output interface{}
		var errorStr *string

		err := rows.Scan(
			&node.ID,
			&node.Type,
			&node.Position.X,
			&node.Position.Y,
			&node.Data.ID,
			&node.Data.Type,
			&node.Data.Label,
			&status,
			&configJSON,
			&outputJSON,
			&errorStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if len(configJSON) > 0 {
			json.Unmarshal(configJSON, &node.Data.Config)
		} else {
			node.Data.Config = make(map[string]interface{})
		}

		if len(outputJSON) > 0 {
			json.Unmarshal(outputJSON, &output)
			node.Data.Output = output
		}

		if status != nil {
			node.Data.Status = *status
		}

		if errorStr != nil {
			node.Data.Error = *errorStr
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// Update updates a node
func (r *NodeRepository) Update(ctx context.Context, flowID uuid.UUID, node *models.FlowNode) error {
	configJSON, _ := json.Marshal(node.Data.Config)
	outputJSON, _ := json.Marshal(node.Data.Output)

	query := `
		UPDATE flow_nodes
		SET type = $3, position_x = $4, position_y = $5,
		    data_type = $6, data_label = $7, data_status = $8,
		    data_config = $9, data_output = $10, data_error = $11,
		    updated_at = $12
		WHERE flow_id = $1 AND node_id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		flowID,
		node.ID,
		node.Type,
		node.Position.X,
		node.Position.Y,
		node.Data.Type,
		node.Data.Label,
		node.Data.Status,
		configJSON,
		outputJSON,
		node.Data.Error,
		time.Now(),
	)

	return err
}

// Delete deletes a node
func (r *NodeRepository) Delete(ctx context.Context, flowID uuid.UUID, nodeID string) error {
	query := `DELETE FROM flow_nodes WHERE flow_id = $1 AND node_id = $2`
	_, err := r.db.Exec(ctx, query, flowID, nodeID)
	return err
}

// DeleteByFlowID deletes all nodes for a flow
func (r *NodeRepository) DeleteByFlowID(ctx context.Context, flowID uuid.UUID) error {
	query := `DELETE FROM flow_nodes WHERE flow_id = $1`
	_, err := r.db.Exec(ctx, query, flowID)
	return err
}

// BulkUpsert creates or updates multiple nodes for a flow
func (r *NodeRepository) BulkUpsert(ctx context.Context, flowID uuid.UUID, nodes []models.FlowNode) error {
	// Delete existing nodes
	if err := r.DeleteByFlowID(ctx, flowID); err != nil {
		return err
	}

	// Insert new nodes
	for _, node := range nodes {
		if err := r.Create(ctx, flowID, &node); err != nil {
			return err
		}
	}

	return nil
}

