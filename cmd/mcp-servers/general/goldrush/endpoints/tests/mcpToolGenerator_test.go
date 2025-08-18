package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGenerateEndpointTools(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test that all endpoint tools are generated
	t.Run("Test GenerateEndpointTools returns all tools", func(t *testing.T) {
		tools := ge.GenerateEndpointTools()

		// Should have 11 tools total
		expectedToolCount := 11
		if len(tools) != expectedToolCount {
			t.Errorf("Expected %d tools, got %d", expectedToolCount, len(tools))
		}

		// Check that each tool has both Tool and Handler
		for i, toolInfo := range tools {
			if toolInfo.Tool.Name == "" {
				t.Errorf("Tool %d has empty name", i)
			}

			if toolInfo.Handler == nil {
				t.Errorf("Tool %d (%s) has nil handler", i, toolInfo.Tool.Name)
			}
		}
	})
}

func TestAllToolsAreGenerated(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tools := ge.GenerateEndpointTools()

	// Define expected tool names
	expectedTools := []string{
		"get_activity_across_all_chains",
		"get_bitcoin_balances_for_HD_address",
		"get_erc20_token_transfer_for_address",
		"get_gas_prices",
		"get_historical_token_prices",
		"get_multichain_balances",
		"get_multichain_transactions",
		"get_native_token_balance_for_address",
		"get_nfts_for_address",
		"get_token_balance_for_address",
		"get_transaction",
	}

	t.Run("Test all expected tools are present", func(t *testing.T) {
		// Create a map of found tools
		foundTools := make(map[string]bool)
		for _, toolInfo := range tools {
			foundTools[toolInfo.Tool.Name] = true
		}

		// Check that all expected tools are found
		for _, expectedTool := range expectedTools {
			if !foundTools[expectedTool] {
				t.Errorf("Expected tool '%s' not found", expectedTool)
			}
		}

		// Check that no unexpected tools are present
		if len(foundTools) != len(expectedTools) {
			t.Errorf("Found %d tools, expected %d", len(foundTools), len(expectedTools))
		}
	})
}

func TestToolDescriptions(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tools := ge.GenerateEndpointTools()

	t.Run("Test all tools have descriptions", func(t *testing.T) {
		for i, toolInfo := range tools {
			if toolInfo.Tool.Description == "" {
				t.Errorf("Tool %d (%s) has empty description", i, toolInfo.Tool.Name)
			}
		}
	})
}

func TestToolInputSchemas(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tools := ge.GenerateEndpointTools()

	t.Run("Test all tools have input schemas", func(t *testing.T) {
		for i, toolInfo := range tools {
			if len(toolInfo.Tool.InputSchema.Properties) == 0 {
				t.Errorf("Tool %d (%s) has no input schema properties", i, toolInfo.Tool.Name)
			}

			if len(toolInfo.Tool.InputSchema.Required) == 0 {
				t.Errorf("Tool %d (%s) has no required parameters", i, toolInfo.Tool.Name)
			}
		}
	})
}
