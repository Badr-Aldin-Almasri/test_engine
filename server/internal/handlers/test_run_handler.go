package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visual-api-testing-platform/server/internal/repository"
	"github.com/visual-api-testing-platform/server/internal/engine"
)

// TestRunHandler handles test run-related HTTP requests
type TestRunHandler struct {
	testRunRepo *repository.TestRunRepository
	flowRepo    *repository.FlowRepository
	flowRunner  *engine.FlowRunner
}

// NewTestRunHandler creates a new test run handler
func NewTestRunHandler(
	testRunRepo *repository.TestRunRepository,
	flowRepo *repository.FlowRepository,
	flowRunner *engine.FlowRunner,
) *TestRunHandler {
	return &TestRunHandler{
		testRunRepo: testRunRepo,
		flowRepo:    flowRepo,
		flowRunner:  flowRunner,
	}
}

// RunFlow handles POST /api/flows/:id/run
func (h *TestRunHandler) RunFlow(c *gin.Context) {
	flowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	flow, err := h.flowRepo.GetByID(c.Request.Context(), flowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow not found"})
		return
	}

	// Execute flow in a goroutine
	go func() {
		ctx := context.Background()
		testRun, err := h.flowRunner.ExecuteFlow(ctx, flow)
		if err != nil {
			return
		}

		// Save test run to database
		h.testRunRepo.Create(ctx, testRun)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Flow execution started",
		"flow_id": flowID,
	})
}

// GetTestRun handles GET /api/test-runs/:id
func (h *TestRunHandler) GetTestRun(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	testRun, err := h.testRunRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test run not found"})
		return
	}

	c.JSON(http.StatusOK, testRun)
}

// GetTestRunsByFlow handles GET /api/flows/:id/test-runs
func (h *TestRunHandler) GetTestRunsByFlow(c *gin.Context) {
	flowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	testRuns, err := h.testRunRepo.GetByFlowID(c.Request.Context(), flowID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, testRuns)
}

