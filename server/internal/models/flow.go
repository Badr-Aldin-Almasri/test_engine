package models

import (
	"time"

	"github.com/google/uuid"
)

// Flow represents a test flow definition
type Flow struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	UserID      uuid.UUID              `json:"user_id" db:"user_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Tags        []string               `json:"tags" db:"tags"`
	Nodes       []FlowNode             `json:"nodes" db:"nodes"`
	Edges       []FlowEdge             `json:"edges" db:"edges"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// FlowNode represents a node in a flow
type FlowNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position Position               `json:"position"`
	Data     NodeData               `json:"data"`
}

// Position represents node position in the editor
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// FlowEdge represents a connection between nodes
type FlowEdge struct {
	ID          string  `json:"id"`
	Source      string  `json:"source"`
	Target      string  `json:"target"`
	SourceHandle *string `json:"sourceHandle,omitempty"`
	TargetHandle *string `json:"targetHandle,omitempty"`
}

// NodeData represents the data stored in a node
type NodeData struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Label  string                 `json:"label"`
	Status string                 `json:"status,omitempty"`
	Config map[string]interface{} `json:"config"`
	Output interface{}            `json:"output,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

