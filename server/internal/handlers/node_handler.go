package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visual-api-testing-platform/server/internal/models"
	"github.com/visual-api-testing-platform/server/internal/node"
	"github.com/visual-api-testing-platform/server/internal/repository"
)

// NodeHandler handles node-related HTTP requests
type NodeHandler struct {
	flowRepo   *repository.FlowRepository
	nodeFactory *node.NodeFactory
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(flowRepo *repository.FlowRepository) *NodeHandler {
	return &NodeHandler{
		flowRepo:    flowRepo,
		nodeFactory: node.NewNodeFactory(),
	}
}

// ExecuteNode handles POST /api/nodes/:flowId/:nodeId/execute
func (h *NodeHandler) ExecuteNode(c *gin.Context) {
	flowID, err := uuid.Parse(c.Param("flowId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	nodeID := c.Param("nodeId")
	if nodeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the flow to verify ownership
	flow, err := h.flowRepo.GetByID(c.Request.Context(), flowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow not found"})
		return
	}

	// Verify the flow belongs to the user
	if flow.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to execute this node"})
		return
	}

	// Find the node in the flow
	var targetNode *models.FlowNode
	for i := range flow.Nodes {
		if flow.Nodes[i].ID == nodeID {
			targetNode = &flow.Nodes[i]
			break
		}
	}

	if targetNode == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Node not found in flow"})
		return
	}

	// Only execute API nodes
	if targetNode.Type != "api" && targetNode.Type != "API" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only API nodes can be executed individually"})
		return
	}

	// Get the latest node config from the request body (if provided)
	// Otherwise use the stored config
	var reqBody struct {
		Config map[string]interface{} `json:"config"`
	}
	
	config := targetNode.Data.Config
	if config == nil {
		config = make(map[string]interface{})
	}
	
	// If request body has config, use it (for testing with latest config before saving)
	if c.ShouldBindJSON(&reqBody) == nil && reqBody.Config != nil {
		// Merge request config with stored config (request takes precedence)
		for k, v := range reqBody.Config {
			config[k] = v
		}
	}

	// Create API node instance
	apiNode, err := h.nodeFactory.CreateNode(
		"api",
		targetNode.ID,
		targetNode.Data.Label,
		config,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Execute the node with a timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	output, err := apiNode.Execute(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"output": nil,
		})
		return
	}

	// Return the output
	c.JSON(http.StatusOK, gin.H{
		"output": output,
		"status": "success",
	})
}

