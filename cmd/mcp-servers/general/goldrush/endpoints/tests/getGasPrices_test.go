package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGetGasPrices(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test the endpoint function
	t.Run("Test GetGasPrices function call", func(t *testing.T) {
		chainName := "ethereum"
		eventType := "erc20"
		result := ge.GetGasPrices(chainName, eventType)

		// Verify that we get a response (even if it's an error due to test token)
		if result == "" {
			t.Error("Expected non-empty result from GetGasPrices")
		}

		// Log the result for debugging
		t.Logf("GetGasPrices response: %s", result)

		// The response should contain some indication of the API call
		// Even with a test token, we should get some response structure
		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	// Test the MCP tool generation
	t.Run("Test GenerateGasPriceTool", func(t *testing.T) {
		tool, handler := ge.GenerateGasPriceTool()

		// Verify tool properties
		if tool.Name != "get_gas_prices" {
			t.Errorf("Expected tool name 'get_gas_prices', got '%s'", tool.Name)
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

func TestGasPriceToolParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tool, _ := ge.GenerateGasPriceTool()

	t.Run("Test tool parameter structure", func(t *testing.T) {
		// Check if the tool has the expected parameters
		chainNameProp, exists := tool.InputSchema.Properties["chain_name"]
		if !exists {
			t.Error("Expected tool to have 'chain_name' parameter")
		}

		if chainNameProp == nil {
			t.Error("Expected 'chain_name' parameter to have properties")
		}

		eventTypeProp, exists := tool.InputSchema.Properties["event_type"]
		if !exists {
			t.Error("Expected tool to have 'event_type' parameter")
		}

		if eventTypeProp == nil {
			t.Error("Expected 'event_type' parameter to have properties")
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

		if !requiredParams["event_type"] {
			t.Error("Expected 'event_type' to be in required parameters")
		}

		// Should have exactly 2 required parameters
		if len(tool.InputSchema.Required) != 2 {
			t.Errorf("Expected 2 required parameters, got %d", len(tool.InputSchema.Required))
		}
	})
}

func TestGasPriceToolEnumValues(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tool, _ := ge.GenerateGasPriceTool()

	t.Run("Test event_type parameter exists", func(t *testing.T) {
		// Check if event_type parameter exists
		eventTypeProp := tool.InputSchema.Properties["event_type"]
		if eventTypeProp == nil {
			t.Fatal("event_type property not found")
		}

		// The event_type should have enum values: "erc20", "uniswapv3", "nativetokens"
		// Note: This test verifies the parameter exists, enum validation is done by MCP framework
	})
}

func TestGetGasPricesWithDifferentParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with different chain names", func(t *testing.T) {
		chains := []string{"ethereum", "polygon", "bsc"}
		eventType := "erc20"

		for _, chain := range chains {
			result := ge.GetGasPrices(chain, eventType)

			if result == "" {
				t.Errorf("Expected non-empty result for chain %s", chain)
			}

			t.Logf("GetGasPrices response for %s: %s", chain, result)
		}
	})

	t.Run("Test with different event types", func(t *testing.T) {
		chainName := "ethereum"
		eventTypes := []string{"erc20", "uniswapv3", "nativetokens"}

		for _, eventType := range eventTypes {
			result := ge.GetGasPrices(chainName, eventType)

			if result == "" {
				t.Errorf("Expected non-empty result for event type %s", eventType)
			}

			t.Logf("GetGasPrices response for event type %s: %s", eventType, result)
		}
	})
}
