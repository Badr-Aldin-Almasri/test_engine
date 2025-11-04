package models

import (
	"time"

	"github.com/google/uuid"
)

// TestRun represents a test execution record
type TestRun struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	FlowID      uuid.UUID              `json:"flow_id" db:"flow_id"`
	FlowName    string                 `json:"flow_name,omitempty"`
	Status      ExecutionStatus        `json:"status" db:"status"`
	StartedAt   time.Time              `json:"started_at" db:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	DurationMs  *int                   `json:"duration_ms,omitempty" db:"duration_ms"`
	NodeResults map[string]NodeResult  `json:"node_results" db:"node_results"`
	Error       string                 `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// NodeResult represents the result of a single node execution
type NodeResult struct {
	Status   ExecutionStatus `json:"status"`
	Output   interface{}    `json:"output,omitempty"`
	Error    string         `json:"error,omitempty"`
	Duration int            `json:"duration"` // in milliseconds
}

// ExecutionStatus represents the status of an execution
type ExecutionStatus string

const (
	ExecutionStatusPending ExecutionStatus = "pending"
	ExecutionStatusRunning ExecutionStatus = "running"
	ExecutionStatusSuccess ExecutionStatus = "success"
	ExecutionStatusFailed  ExecutionStatus = "failed"
	ExecutionStatusSkipped ExecutionStatus = "skipped"
)

