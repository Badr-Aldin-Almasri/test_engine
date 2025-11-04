package node

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// VerificationNode performs assertions on data
type VerificationNode struct {
	BaseNode
}

// NewVerificationNode creates a new Verification node
func NewVerificationNode(id, label string, config map[string]interface{}) *VerificationNode {
	return &VerificationNode{
		BaseNode: BaseNode{
			ID:     id,
			Label:  label,
			Config: config,
		},
	}
}

// GetType returns the node type
func (v *VerificationNode) GetType() NodeType {
	return NodeTypeVerification
}

// ValidateConfig validates the Verification node configuration
func (v *VerificationNode) ValidateConfig() error {
	assertionType, ok := v.Config["assertionType"].(string)
	if !ok {
		assertionType = "equals" // default
	}

	validTypes := map[string]bool{
		"equals": true, "contains": true, "regex": true, "custom": true,
	}
	if !validTypes[assertionType] {
		return fmt.Errorf("invalid assertion type: %s", assertionType)
	}

	return nil
}

// Execute performs the verification
func (v *VerificationNode) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if err := v.ValidateConfig(); err != nil {
		return nil, err
	}

	assertionType, _ := v.Config["assertionType"].(string)
	if assertionType == "" {
		assertionType = "equals"
	}

	expected, exists := v.Config["expected"]
	if !exists {
		return nil, fmt.Errorf("expected value is required")
	}

	// Get actual value from input (usually from previous node)
	actual, ok := input["data"].(map[string]interface{})
	if !ok {
		actual = input // Use entire input as actual
	}

	var passed bool
	var err error

	switch assertionType {
	case "equals":
		passed, err = v.equalsAssertion(expected, actual)
	case "contains":
		passed, err = v.containsAssertion(expected, actual)
	case "regex":
		passed, err = v.regexAssertion(expected, actual)
	case "custom":
		// Custom script execution would go here
		// For now, we'll return a placeholder
		passed = true // Placeholder
	default:
		return nil, fmt.Errorf("unsupported assertion type: %s", assertionType)
	}

	if err != nil {
		return nil, err
	}

	if !passed {
		return nil, fmt.Errorf("verification failed: expected %v, got %v", expected, actual)
	}

	return map[string]interface{}{
		"passed":  true,
		"expected": expected,
		"actual":   actual,
	}, nil
}

func (v *VerificationNode) equalsAssertion(expected, actual interface{}) (bool, error) {
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		return false, err
	}
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		return false, err
	}
	return string(expectedJSON) == string(actualJSON), nil
}

func (v *VerificationNode) containsAssertion(expected, actual interface{}) (bool, error) {
	expectedStr := fmt.Sprintf("%v", expected)
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(actualJSON), expectedStr), nil
}

func (v *VerificationNode) regexAssertion(expected, actual interface{}) (bool, error) {
	pattern, ok := expected.(string)
	if !ok {
		return false, fmt.Errorf("regex pattern must be a string")
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		return false, err
	}

	matched, err := regexp.MatchString(pattern, string(actualJSON))
	return matched, err
}

