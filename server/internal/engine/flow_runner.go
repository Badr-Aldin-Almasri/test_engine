package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/visual-api-testing-platform/server/internal/models"
	"github.com/visual-api-testing-platform/server/internal/node"
)

// FlowRunner executes flows concurrently
type FlowRunner struct {
	nodeFactory *node.NodeFactory
	hub         *ExecutionHub
}

// NewFlowRunner creates a new flow runner
func NewFlowRunner(hub *ExecutionHub) *FlowRunner {
	return &FlowRunner{
		nodeFactory: node.NewNodeFactory(),
		hub:         hub,
	}
}

// ExecuteFlow executes a flow and returns the test run result
func (r *FlowRunner) ExecuteFlow(ctx context.Context, flow *models.Flow) (*models.TestRun, error) {
	testRun := &models.TestRun{
		ID:          uuid.New(),
		FlowID:      flow.ID,
		FlowName:    flow.Name,
		Status:      models.ExecutionStatusRunning,
		StartedAt:   time.Now(),
		NodeResults: make(map[string]models.NodeResult),
	}

	// Build dependency graph
	nodeMap := make(map[string]*models.FlowNode)
	incomingEdges := make(map[string][]models.FlowEdge)

	for i := range flow.Nodes {
		nodeMap[flow.Nodes[i].ID] = &flow.Nodes[i]
	}

	for _, edge := range flow.Edges {
		if incomingEdges[edge.Target] == nil {
			incomingEdges[edge.Target] = make([]models.FlowEdge, 0)
		}
		incomingEdges[edge.Target] = append(incomingEdges[edge.Target], edge)
	}

	// Find root nodes (nodes with no incoming edges)
	rootNodes := make([]*models.FlowNode, 0)
	for _, n := range flow.Nodes {
		if len(incomingEdges[n.ID]) == 0 {
			rootNodes = append(rootNodes, &n)
		}
	}

	// Execute nodes concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	executedNodes := make(map[string]bool)
	nodeOutputs := make(map[string]map[string]interface{})

	// Function to execute a node (declared first for recursive calls)
	var executeNode func(*models.FlowNode)
	executeNode = func(n *models.FlowNode) {
		defer wg.Done()

		// Check if dependencies are ready
		deps := incomingEdges[n.ID]
		for _, dep := range deps {
			mu.Lock()
			if !executedNodes[dep.Source] {
				mu.Unlock()
				// Dependency not ready, retry later
				time.Sleep(100 * time.Millisecond)
				wg.Add(1)
				go executeNode(n)
				return
			}
			mu.Unlock()
		}

		// Build input from dependencies
		input := make(map[string]interface{})
		for _, dep := range deps {
			if output, ok := nodeOutputs[dep.Source]; ok {
				input[dep.Source] = output
				// Also merge into top-level data
				if data, ok := output["data"].(map[string]interface{}); ok {
					for k, v := range data {
						input[k] = v
					}
				}
			}
		}

		// Create and execute node
		startTime := time.Now()
		r.hub.BroadcastNodeUpdate(testRun.ID, n.ID, "running", nil, "")

		nodeInstance, err := r.nodeFactory.CreateNode(
			n.Data.Type,
			n.Data.ID,
			n.Data.Label,
			n.Data.Config,
		)
		if err != nil {
			mu.Lock()
			testRun.NodeResults[n.ID] = models.NodeResult{
				Status:   models.ExecutionStatusFailed,
				Error:    err.Error(),
				Duration: int(time.Since(startTime).Milliseconds()),
			}
			mu.Unlock()
			r.hub.BroadcastNodeUpdate(testRun.ID, n.ID, "failed", nil, err.Error())
			return
		}

		output, err := nodeInstance.Execute(ctx, input)
		duration := int(time.Since(startTime).Milliseconds())

		mu.Lock()
		executedNodes[n.ID] = true
		if err != nil {
			testRun.NodeResults[n.ID] = models.NodeResult{
				Status:   models.ExecutionStatusFailed,
				Error:    err.Error(),
				Duration: duration,
			}
			r.hub.BroadcastNodeUpdate(testRun.ID, n.ID, "failed", nil, err.Error())
		} else {
			nodeOutputs[n.ID] = output
			testRun.NodeResults[n.ID] = models.NodeResult{
				Status:   models.ExecutionStatusSuccess,
				Output:   output,
				Duration: duration,
			}
			r.hub.BroadcastNodeUpdate(testRun.ID, n.ID, "success", output, "")
		}
		mu.Unlock()
	}

	// Execute root nodes first
	for _, rootNode := range rootNodes {
		wg.Add(1)
		go executeNode(rootNode)
	}

	// Execute remaining nodes
	for _, n := range flow.Nodes {
		if len(incomingEdges[n.ID]) > 0 {
			wg.Add(1)
			go executeNode(&n)
		}
	}

	// Wait for all nodes to complete (with timeout)
	done := make(chan bool)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// All nodes completed
	case <-ctx.Done():
		testRun.Status = models.ExecutionStatusFailed
		testRun.Error = "Execution cancelled"
		return testRun, ctx.Err()
	case <-time.After(5 * time.Minute):
		testRun.Status = models.ExecutionStatusFailed
		testRun.Error = "Execution timeout"
		return testRun, fmt.Errorf("execution timeout")
	}

	// Check if any node failed
	for _, result := range testRun.NodeResults {
		if result.Status == models.ExecutionStatusFailed {
			testRun.Status = models.ExecutionStatusFailed
			break
		}
	}

	if testRun.Status != models.ExecutionStatusFailed {
		testRun.Status = models.ExecutionStatusSuccess
	}

	completedAt := time.Now()
	testRun.CompletedAt = &completedAt
	duration := int(time.Since(testRun.StartedAt).Milliseconds())
	testRun.DurationMs = &duration

	r.hub.BroadcastTestRunComplete(testRun)

	return testRun, nil
}

