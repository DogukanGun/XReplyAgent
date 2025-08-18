package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestGetActivityAcrossAllChains(t *testing.T) {
	// Set up test environment
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	// Create a new GoldrushEndpoints instance
	ge := endpoints.NewGoldrushEndpoints()

	// Test the endpoint function
	t.Run("Test GetActivityAcrossAllChains function call", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// Verify that we get a response (even if it's an error due to test token)
		if result == "" {
			t.Error("Expected non-empty result from GetActivityAcrossAllChains")
		}

		// Log the result for debugging
		t.Logf("GetActivityAcrossAllChains response: %s", result)

		// The response should contain some indication of the API call
		// Even with a test token, we should get some response structure
		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	// Test the MCP tool generation
	t.Run("Test GenerateActivityTool", func(t *testing.T) {
		tool, handler := ge.GenerateActivityTool()

		// Verify tool properties
		if tool.Name != "get_activity_across_all_chains" {
			t.Errorf("Expected tool name 'get_activity_across_all_chains', got '%s'", tool.Name)
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

func TestGetActivityAcrossAllChainsWithInvalidAddress(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with empty wallet address", func(t *testing.T) {
		result := ge.GetActivityAcrossAllChains("")

		// Should handle empty address gracefully and still return a response
		if result == "" {
			t.Error("Expected result even with empty wallet address")
		}

		t.Logf("Empty address response: %s", result)
	})

	t.Run("Test with invalid wallet address format", func(t *testing.T) {
		result := ge.GetActivityAcrossAllChains("invalid_address")

		// Should handle invalid address format gracefully and return a response
		if result == "" {
			t.Error("Expected result even with invalid wallet address format")
		}

		t.Logf("Invalid address response: %s", result)
	})
}

func TestGetActivityAcrossAllChainsWithRealToken(t *testing.T) {
	// Test with actual auth token if available
	authToken := os.Getenv("GOLDRUSH_AUTH_TOKEN")
	if authToken == "" || authToken == "test_token" {
		t.Skip("Skipping real token test - no valid GOLDRUSH_AUTH_TOKEN found")
	}

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with real auth token", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// With a real token, we should get a proper API response
		if result == "" {
			t.Error("Expected non-empty result with real auth token")
		}

		t.Logf("Real token response length: %d", len(result))

		// The response should be substantial with real data
		if len(result) < 100 {
			t.Errorf("Expected substantial response with real token, got %d characters", len(result))
		}
	})
}

func TestGoldrushEndpointsInitialization(t *testing.T) {
	t.Run("Test NewGoldrushEndpoints", func(t *testing.T) {
		os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
		defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

		ge := endpoints.NewGoldrushEndpoints()

		if ge == nil {
			t.Error("Expected non-nil GoldrushEndpoints instance")
		}

		if ge.BaseUrl != "https://api.covalenthq.com/v1/" {
			t.Errorf("Expected BaseUrl 'https://api.covalenthq.com/v1/', got '%s'", ge.BaseUrl)
		}

		if ge.AuthToken != "test_token" {
			t.Errorf("Expected AuthToken 'test_token', got '%s'", ge.AuthToken)
		}
	})

	t.Run("Test NewGoldrushEndpoints without auth token", func(t *testing.T) {
		os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

		ge := endpoints.NewGoldrushEndpoints()

		if ge == nil {
			t.Error("Expected non-nil GoldrushEndpoints instance even without auth token")
		}

		if ge.AuthToken != "" {
			t.Errorf("Expected empty AuthToken when env var not set, got '%s'", ge.AuthToken)
		}
	})
}
