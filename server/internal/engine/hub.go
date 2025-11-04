package engine

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/visual-api-testing-platform/server/internal/models"
)

// ExecutionHub manages WebSocket connections for real-time updates
type ExecutionHub struct {
	clients    map[uuid.UUID]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// Client represents a WebSocket client
type Client struct {
	hub       *ExecutionHub
	send      chan []byte
	testRunID uuid.UUID
}

// NewClient creates a new WebSocket client
func NewClient(hub *ExecutionHub, testRunID uuid.UUID) *Client {
	return &Client{
		hub:       hub,
		send:      make(chan []byte, 256),
		testRunID: testRunID,
	}
}

// NewExecutionHub creates a new execution hub
func NewExecutionHub() *ExecutionHub {
	return &ExecutionHub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *ExecutionHub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.testRunID] == nil {
				h.clients[client.testRunID] = make(map[*Client]bool)
			}
			h.clients[client.testRunID][client] = true

		case client := <-h.unregister:
			if clients, ok := h.clients[client.testRunID]; ok {
				delete(clients, client)
				close(client.send)
				if len(clients) == 0 {
					delete(h.clients, client.testRunID)
				}
			}

		case message := <-h.broadcast:
			// Parse message to get testRunID
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if testRunIDStr, ok := msg["testRunId"].(string); ok {
					if testRunID, err := uuid.Parse(testRunIDStr); err == nil {
						if clients, ok := h.clients[testRunID]; ok {
							for client := range clients {
								select {
								case client.send <- message:
								default:
									close(client.send)
									delete(clients, client)
								}
							}
						}
					}
				}
			}
		}
	}
}

// BroadcastNodeUpdate broadcasts a node execution update
func (h *ExecutionHub) BroadcastNodeUpdate(testRunID uuid.UUID, nodeID, status string, output interface{}, error string) {
	message := map[string]interface{}{
		"type":      "node_update",
		"testRunId": testRunID.String(),
		"nodeId":    nodeID,
		"status":    status,
		"output":    output,
		"error":     error,
	}

	data, _ := json.Marshal(message)
	h.broadcast <- data
}

// BroadcastTestRunComplete broadcasts test run completion
func (h *ExecutionHub) BroadcastTestRunComplete(testRun *models.TestRun) {
	message := map[string]interface{}{
		"type":      "test_run_complete",
		"testRunId": testRun.ID.String(),
		"status":    string(testRun.Status),
		"duration":  testRun.DurationMs,
	}

	data, _ := json.Marshal(message)
	h.broadcast <- data
}

// Register registers a new client
func (h *ExecutionHub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client
func (h *ExecutionHub) Unregister(client *Client) {
	h.unregister <- client
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		c.hub.Unregister(c)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump(conn *websocket.Conn) {
	defer func() {
		conn.Close()
		c.hub.Unregister(c)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

