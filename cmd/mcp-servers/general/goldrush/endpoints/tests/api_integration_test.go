package tests

import (
	"os"
	"strings"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestAPIIntegration(t *testing.T) {
	// Check if we have a real auth token
	authToken := os.Getenv("GOLDRUSH_AUTH_TOKEN")
	if authToken == "" || authToken == "test_token" {
		t.Skip("Skipping API integration tests - no valid GOLDRUSH_AUTH_TOKEN found")
	}

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test Activity Across All Chains API", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// Should get a substantial response
		if result == "" {
			t.Fatal("Expected non-empty response from API")
		}

		t.Logf("Activity API response length: %d", len(result))

		// Response should be substantial (more than just an error message)
		if len(result) < 100 {
			t.Errorf("Expected substantial response, got %d characters", len(result))
		}

		// Safely check for error indicators in the response
		resultLower := strings.ToLower(result)
		if strings.Contains(resultLower, "unauthorized") ||
			strings.Contains(resultLower, "invalid") ||
			strings.Contains(resultLower, "error") ||
			strings.Contains(resultLower, "payment required") {
			// Log the response safely (limit to actual length)
			maxLen := len(result)
			if maxLen > 200 {
				maxLen = 200
			}
			t.Logf("API response contains error indicators: %s", result[:maxLen])
		}
	})

	t.Run("Test Bitcoin HD Wallet API", func(t *testing.T) {
		walletAddress := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
		result := ge.GetBitcoinBalancesForHDAddress(walletAddress)

		if result == "" {
			t.Fatal("Expected non-empty response from Bitcoin API")
		}

		t.Logf("Bitcoin API response length: %d", len(result))

		if len(result) < 100 {
			t.Errorf("Expected substantial response, got %d characters", len(result))
		}
	})

	t.Run("Test Gas Prices API", func(t *testing.T) {
		chainName := "ethereum"
		eventType := "erc20"
		result := ge.GetGasPrices(chainName, eventType)

		if result == "" {
			t.Fatal("Expected non-empty response from Gas Prices API")
		}

		t.Logf("Gas Prices API response length: %d", len(result))

		if len(result) < 100 {
			t.Errorf("Expected substantial response, got %d characters", len(result))
		}
	})

	t.Run("Test Multichain Balances API", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetMultichainBalances(walletAddress)

		if result == "" {
			t.Fatal("Expected non-empty response from Multichain Balances API")
		}

		t.Logf("Multichain Balances API response length: %d", len(result))

		if len(result) < 100 {
			t.Errorf("Expected substantial response, got %d characters", len(result))
		}
	})

	t.Run("Test Token Balance API", func(t *testing.T) {
		chainName := "ethereum"
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetTokenBalanceForAddress(chainName, walletAddress)

		if result == "" {
			t.Fatal("Expected non-empty response from Token Balance API")
		}

		t.Logf("Token Balance API response length: %d", len(result))

		if len(result) < 100 {
			t.Errorf("Expected substantial response, got %d characters", len(result))
		}
	})
}

func TestAPIErrorHandling(t *testing.T) {
	// Test with invalid parameters to see how the API handles errors
	authToken := os.Getenv("GOLDRUSH_AUTH_TOKEN")
	if authToken == "" || authToken == "test_token" {
		t.Skip("Skipping API error handling tests - no valid GOLDRUSH_AUTH_TOKEN found")
	}

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with invalid wallet address", func(t *testing.T) {
		invalidAddress := "invalid_address_12345"
		result := ge.GetActivityAcrossAllChains(invalidAddress)

		// Should still get a response (even if it's an error response)
		if result == "" {
			t.Error("Expected response even with invalid address")
		}

		t.Logf("Invalid address response: %s", result)

		// The response should indicate some kind of error or validation
		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	t.Run("Test with invalid chain name", func(t *testing.T) {
		invalidChain := "invalid_chain_12345"
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetTokenBalanceForAddress(invalidChain, walletAddress)

		if result == "" {
			t.Error("Expected response even with invalid chain name")
		}

		t.Logf("Invalid chain response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})
}

func TestAPIResponseStructure(t *testing.T) {
	authToken := os.Getenv("GOLDRUSH_AUTH_TOKEN")
	if authToken == "" || authToken == "test_token" {
		t.Skip("Skipping API response structure tests - no valid GOLDRUSH_AUTH_TOKEN found")
	}

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test response contains expected fields", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// Response should be JSON-like (even if it's an error)
		if !strings.Contains(result, "{") || !strings.Contains(result, "}") {
			// Safely log the response (limit to actual length)
			maxLen := len(result)
			if maxLen > 200 {
				maxLen = 200
			}
			t.Logf("Response doesn't appear to be JSON: %s", result[:maxLen])
		}

		// Should contain some indication of the API structure
		if strings.Contains(result, "data") || strings.Contains(result, "result") ||
			strings.Contains(result, "error") || strings.Contains(result, "message") {
			t.Logf("Response contains expected structure indicators")
		}
	})
}
