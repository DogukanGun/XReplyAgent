package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGetBitcoinBalancesForHDAddress(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test the endpoint function
	t.Run("Test GetBitcoinBalancesForHDAddress function call", func(t *testing.T) {
		walletAddress := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
		result := ge.GetBitcoinBalancesForHDAddress(walletAddress)

		// Verify that we get a response (even if it's an error due to test token)
		if result == "" {
			t.Error("Expected non-empty result from GetBitcoinBalancesForHDAddress")
		}

		// Log the result for debugging
		t.Logf("GetBitcoinBalancesForHDAddress response: %s", result)

		// The response should contain some indication of the API call
		// Even with a test token, we should get some response structure
		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	// Test the MCP tool generation
	t.Run("Test GenerateBitcoinBalanceTool", func(t *testing.T) {
		tool, handler := ge.GenerateBitcoinBalanceTool()

		// Verify tool properties
		if tool.Name != "get_bitcoin_balances_for_HD_address" {
			t.Errorf("Expected tool name 'get_bitcoin_balances_for_HD_address', got '%s'", tool.Name)
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

func TestBitcoinBalanceToolParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()
	tool, _ := ge.GenerateBitcoinBalanceTool()

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
	})
}

func TestGetBitcoinBalancesForHDAddressWithInvalidAddress(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with empty wallet address", func(t *testing.T) {
		result := ge.GetBitcoinBalancesForHDAddress("")

		// Should handle empty address gracefully and still return a response
		if result == "" {
			t.Error("Expected result even with empty wallet address")
		}

		t.Logf("Empty address response: %s", result)
	})

	t.Run("Test with invalid wallet address format", func(t *testing.T) {
		result := ge.GetBitcoinBalancesForHDAddress("invalid_address")

		// Should handle invalid address format gracefully and return a response
		if result == "" {
			t.Error("Expected result even with invalid wallet address format")
		}

		t.Logf("Invalid address response: %s", result)
	})
}
