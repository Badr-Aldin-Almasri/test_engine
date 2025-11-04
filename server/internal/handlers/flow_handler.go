package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/visual-api-testing-platform/server/internal/models"
	"github.com/visual-api-testing-platform/server/internal/repository"
)

// FlowHandler handles flow-related HTTP requests
type FlowHandler struct {
	flowRepo *repository.FlowRepository
}

// NewFlowHandler creates a new flow handler
func NewFlowHandler(flowRepo *repository.FlowRepository) *FlowHandler {
	return &FlowHandler{flowRepo: flowRepo}
}

// CreateFlow handles POST /api/flows
func (h *FlowHandler) CreateFlow(c *gin.Context) {
	var req struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		Tags        []string               `json:"tags"`
		Nodes       []models.FlowNode      `json:"nodes"`
		Edges       []models.FlowEdge      `json:"edges"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	flow := &models.Flow{
		ID:          uuid.New(),
		UserID:      userID.(uuid.UUID),
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		Nodes:       req.Nodes,
		Edges:       req.Edges,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.flowRepo.Create(c.Request.Context(), flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flow)
} 

// GetFlow handles GET /api/flows/:id
func (h *FlowHandler) GetFlow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	flow, err := h.flowRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow not found"})
		return
	}

	c.JSON(http.StatusOK, flow)
}

// ListFlows handles GET /api/flows
func (h *FlowHandler) ListFlows(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	flows, err := h.flowRepo.GetByUserID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flows)
}

// UpdateFlow handles PUT /api/flows/:id
func (h *FlowHandler) UpdateFlow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Tags        []string          `json:"tags"`
		Nodes       []models.FlowNode `json:"nodes"`
		Edges       []models.FlowEdge `json:"edges"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flow, err := h.flowRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow not found"})
		return
	}

	// Verify the flow belongs to the user
	if flow.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this flow"})
		return
	}

	if req.Name != "" {
		flow.Name = req.Name
	}
	if req.Description != "" {
		flow.Description = req.Description
	}
	if req.Tags != nil {
		flow.Tags = req.Tags
	}
	if req.Nodes != nil {
		flow.Nodes = req.Nodes
	}
	if req.Edges != nil {
		flow.Edges = req.Edges
	}
	flow.UpdatedAt = time.Now()

	if err := h.flowRepo.Update(c.Request.Context(), flow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flow)
}

// DeleteFlow handles DELETE /api/flows/:id
func (h *FlowHandler) DeleteFlow(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flow ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify the flow belongs to the user
	flow, err := h.flowRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flow not found"})
		return
	}

	if flow.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this flow"})
		return
	}

	if err := h.flowRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flow deleted"})
}

