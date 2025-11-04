package node

import "fmt"

// NodeFactory creates nodes based on their type
type NodeFactory struct{}

// NewNodeFactory creates a new node factory
func NewNodeFactory() *NodeFactory {
	return &NodeFactory{}
}

// CreateNode creates a node instance based on the node type
func (f *NodeFactory) CreateNode(nodeType, id, label string, config map[string]interface{}) (Node, error) {
	nt := NodeType(nodeType)

	switch nt {
	case NodeTypeAPI:
		return NewAPINode(id, label, config), nil
	case NodeTypeMock:
		return NewMockNode(id, label, config), nil
	case NodeTypeVerification:
		return NewVerificationNode(id, label, config), nil
	case NodeTypeReport:
		return NewReportNode(id, label, config), nil
	case NodeTypeEventTrigger:
		// Event trigger would be handled differently in production
		// For now, return a placeholder
		return NewMockNode(id, label, config), nil
	default:
		return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}

