package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGetHistoricalTokenPrices(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test the MCP tool generation
	t.Run("Test GenerateHistoricalTokenPriceTool", func(t *testing.T) {
		tool, handler := ge.GenerateHistoricalTokenPriceTool()

		// Verify tool properties
		if tool.Name != "get_historical_token_prices" {
			t.Errorf("Expected tool name 'get_historical_token_prices', got '%s'", tool.Name)
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

func TestHistoricalTokenPriceToolParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tool, _ := ge.GenerateHistoricalTokenPriceTool()

	t.Run("Test tool parameter structure", func(t *testing.T) {
		// Check if the tool has the expected parameters
		chainNameProp, exists := tool.InputSchema.Properties["chain_name"]
		if !exists {
			t.Error("Expected tool to have 'chain_name' parameter")
		}

		if chainNameProp == nil {
			t.Error("Expected 'chain_name' parameter to have properties")
		}

		contractAddressProp, exists := tool.InputSchema.Properties["contract_address"]
		if !exists {
			t.Error("Expected tool to have 'contract_address' parameter")
		}

		if contractAddressProp == nil {
			t.Error("Expected 'contract_address' parameter to have properties")
		}
	})

	t.Run("Test tool required parameters", func(t *testing.T) {
		// Check if both parameters are in required parameters
		requiredParams := make(map[string]bool)
		for _, required := range tool.InputSchema.Required {
			requiredParams[required] = true
		}

		if !requiredParams["chain_name"] {
			t.Error("Expected 'chain_name' to be in required parameters")
		}

		if !requiredParams["contract_address"] {
			t.Error("Expected 'contract_address' to be in required parameters")
		}

		// Should have exactly 2 required parameters
		if len(tool.InputSchema.Required) != 2 {
			t.Errorf("Expected 2 required parameters, got %d", len(tool.InputSchema.Required))
		}
	})
}
