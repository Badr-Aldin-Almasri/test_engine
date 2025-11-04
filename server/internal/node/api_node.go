package node

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APINode executes HTTP requests
type APINode struct {
	BaseNode
	Client *http.Client
}

// NewAPINode creates a new API node
func NewAPINode(id, label string, config map[string]interface{}) *APINode {
	return &APINode{
		BaseNode: BaseNode{
			ID:     id,
			Label:  label,
			Config: config,
		},
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetType returns the node type
func (a *APINode) GetType() NodeType {
	return NodeTypeAPI
}

// ValidateConfig validates the API node configuration
func (a *APINode) ValidateConfig() error {
	method, ok := a.Config["method"].(string)
	if !ok || method == "" {
		return fmt.Errorf("method is required")
	}

	url, ok := a.Config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("url is required")
	}

	// Validate HTTP method
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
	}
	if !validMethods[method] {
		return fmt.Errorf("invalid HTTP method: %s", method)
	}

	return nil
}

// Execute performs the HTTP request
func (a *APINode) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if err := a.ValidateConfig(); err != nil {
		return nil, err
	}

	method := a.Config["method"].(string)
	url := a.Config["url"].(string)

	// Create request
	var body io.Reader
	if bodyStr, ok := a.Config["body"].(string); ok && bodyStr != "" {
		body = bytes.NewBufferString(bodyStr)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if headers, ok := a.Config["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			if str, ok := v.(string); ok {
				req.Header.Set(k, str)
			}
		}
	}

	// Add default Content-Type if not present and body exists
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response if possible
	var jsonData interface{}
	if err := json.Unmarshal(respBody, &jsonData); err != nil {
		// If not JSON, return as string
		jsonData = string(respBody)
	}

	result := map[string]interface{}{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    resp.Header,
		"data":       jsonData,
	}

	return result, nil
}

