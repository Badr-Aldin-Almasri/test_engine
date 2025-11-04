package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/visual-api-testing-platform/server/internal/engine"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *engine.ExecutionHub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *engine.ExecutionHub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	testRunIDStr := c.Query("testRunId")
	if testRunIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "testRunId query parameter required"})
		return
	}

	testRunID, err := uuid.Parse(testRunIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid testRunId"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := engine.NewClient(h.hub, testRunID)
	h.hub.Register(client)

	go client.WritePump(conn)
	go client.ReadPump(conn)
}

