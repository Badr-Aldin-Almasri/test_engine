package node

import (
	"context"
)

// NodeType represents the type of a node
type NodeType string

const (
	NodeTypeAPI           NodeType = "api"
	NodeTypeVerification NodeType = "verification"
	NodeTypeMock          NodeType = "mock"
	NodeTypeReport        NodeType = "report"
	NodeTypeEventTrigger  NodeType = "event_trigger"
)

// Node is the interface that all node types must implement
type Node interface {
	// Execute runs the node logic and returns output data
	Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
	
	// GetType returns the type of the node
	GetType() NodeType
	
	// ValidateConfig validates the node configuration
	ValidateConfig() error
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	ID     string
	Label  string
	Config map[string]interface{}
}

// GetID returns the node ID
func (b *BaseNode) GetID() string {
	return b.ID
}

// GetLabel returns the node label
func (b *BaseNode) GetLabel() string {
	return b.Label
}

// GetConfig returns the node configuration
func (b *BaseNode) GetConfig() map[string]interface{} {
	return b.Config
}

