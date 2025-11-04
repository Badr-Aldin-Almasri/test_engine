package node

import (
	"context"
	"time"
)

// ReportNode generates test reports
type ReportNode struct {
	BaseNode
}

// NewReportNode creates a new Report node
func NewReportNode(id, label string, config map[string]interface{}) *ReportNode {
	return &ReportNode{
		BaseNode: BaseNode{
			ID:     id,
			Label:  label,
			Config: config,
		},
	}
}

// GetType returns the node type
func (r *ReportNode) GetType() NodeType {
	return NodeTypeReport
}

// ValidateConfig validates the Report node configuration
func (r *ReportNode) ValidateConfig() error {
	// Report node doesn't require strict validation
	return nil
}

// Execute generates a report
func (r *ReportNode) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	reportName, _ := r.Config["reportName"].(string)
	if reportName == "" {
		reportName = "Test Report"
	}

	report := map[string]interface{}{
		"name":      reportName,
		"timestamp": time.Now().Format(time.RFC3339),
		"input":     input,
		"status":    "completed",
	}

	return report, nil
}

