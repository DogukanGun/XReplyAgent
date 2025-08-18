package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGetMultichainTransactions(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test the MCP tool generation
	t.Run("Test GenerateMultichainTransactionsTool", func(t *testing.T) {
		tool, handler := ge.GenerateMultichainTransactionsTool()

		// Verify tool properties
		if tool.Name != "get_multichain_transactions" {
			t.Errorf("Expected tool name 'get_multichain_transactions', got '%s'", tool.Name)
		}

		if tool.Description == "" {
			t.Error("Expected non-empty tool description")
		}

		// Verify tool has required parameters
		if len(tool.InputSchema.Properties) == 0 {
			t.Error("Expected tool to have input schema properties")
		}

		// Verify handler exists
		if handler == nil {
			t.Error("Expected non-nil handler")
		}
	})
}

func TestMultichainTransactionsToolParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tool, _ := ge.GenerateMultichainTransactionsTool()

	t.Run("Test tool parameter structure", func(t *testing.T) {
		// Check if the tool has the expected parameter
		walletAddressProp, exists := tool.InputSchema.Properties["wallet_address"]
		if !exists {
			t.Error("Expected tool to have 'wallet_address' parameter")
		}

		if walletAddressProp == nil {
			t.Error("Expected 'wallet_address' parameter to have properties")
		}
	})

	t.Run("Test tool required parameters", func(t *testing.T) {
		// Check if wallet_address is in required parameters
		found := false
		for _, required := range tool.InputSchema.Required {
			if required == "wallet_address" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected 'wallet_address' to be in required parameters")
		}

		// Should have exactly 1 required parameter
		if len(tool.InputSchema.Required) != 1 {
			t.Errorf("Expected 1 required parameter, got %d", len(tool.InputSchema.Required))
		}
	})
}
