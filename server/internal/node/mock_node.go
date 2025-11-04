package node

import (
	"context"
	"time"
)

// MockNode returns predefined mock responses
type MockNode struct {
	BaseNode
}

// NewMockNode creates a new Mock node
func NewMockNode(id, label string, config map[string]interface{}) *MockNode {
	return &MockNode{
		BaseNode: BaseNode{
			ID:     id,
			Label:  label,
			Config: config,
		},
	}
}

// GetType returns the node type
func (m *MockNode) GetType() NodeType {
	return NodeTypeMock
}

// ValidateConfig validates the Mock node configuration
func (m *MockNode) ValidateConfig() error {
	// Mock node doesn't require strict validation
	// If mockResponse is not set, we'll use a default
	return nil
}

// Execute returns the mock response
func (m *MockNode) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	// Simulate delay if configured
	if delay, ok := m.Config["mockDelay"].(float64); ok && delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(delay) * time.Millisecond):
		}
	}

	// Return mock response or default
	if mockResponse, ok := m.Config["mockResponse"]; ok && mockResponse != nil {
		if result, ok := mockResponse.(map[string]interface{}); ok {
			return result, nil
		}
		// If it's not a map, wrap it
		return map[string]interface{}{
			"status": 200,
			"data":   mockResponse,
		}, nil
	}

	// Default mock response
	return map[string]interface{}{
		"status": 200,
		"data": map[string]interface{}{
			"message":   "Mock response",
			"timestamp": time.Now().Unix(),
		},
	}, nil
}

