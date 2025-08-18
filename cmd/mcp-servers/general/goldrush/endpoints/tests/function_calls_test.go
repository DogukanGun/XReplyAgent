package tests

import (
	"os"
	"testing"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"
)

func TestFunctionCallsWithTestToken(t *testing.T) {
	// Set up test environment with test token
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test GetActivityAcrossAllChains function call", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// Even with test token, we should get some response
		if result == "" {
			t.Error("Expected non-empty result from GetActivityAcrossAllChains")
		}

		t.Logf("GetActivityAcrossAllChains response: %s", result)

		// Response should be at least 10 characters (even if it's an error message)
		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	t.Run("Test GetBitcoinBalancesForHDAddress function call", func(t *testing.T) {
		walletAddress := "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"
		result := ge.GetBitcoinBalancesForHDAddress(walletAddress)

		if result == "" {
			t.Error("Expected non-empty result from GetBitcoinBalancesForHDAddress")
		}

		t.Logf("GetBitcoinBalancesForHDAddress response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	t.Run("Test GetGasPrices function call", func(t *testing.T) {
		chainName := "ethereum"
		eventType := "erc20"
		result := ge.GetGasPrices(chainName, eventType)

		if result == "" {
			t.Error("Expected non-empty result from GetGasPrices")
		}

		t.Logf("GetGasPrices response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	t.Run("Test GetMultichainBalances function call", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetMultichainBalances(walletAddress)

		if result == "" {
			t.Error("Expected non-empty result from GetMultichainBalances")
		}

		t.Logf("GetMultichainBalances response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})

	t.Run("Test GetTokenBalanceForAddress function call", func(t *testing.T) {
		chainName := "ethereum"
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetTokenBalanceForAddress(chainName, walletAddress)

		if result == "" {
			t.Error("Expected non-empty result from GetTokenBalanceForAddress")
		}

		t.Logf("GetTokenBalanceForAddress response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})
}

func TestFunctionCallsWithEmptyToken(t *testing.T) {
	// Test with no auth token
	os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test GetActivityAcrossAllChains with no token", func(t *testing.T) {
		walletAddress := "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6"
		result := ge.GetActivityAcrossAllChains(walletAddress)

		// Should still get a response (even if it's an error)
		if result == "" {
			t.Error("Expected result even with no auth token")
		}

		t.Logf("No token response: %s", result)

		if len(result) < 10 {
			t.Errorf("Expected response to be longer than 10 characters, got %d", len(result))
		}
	})
}

func TestFunctionCallsWithInvalidParameters(t *testing.T) {
	os.Setenv("GOLDRUSH_AUTH_TOKEN", "test_token")
	defer os.Unsetenv("GOLDRUSH_AUTH_TOKEN")

	ge := endpoints.NewGoldrushEndpoints()

	t.Run("Test with empty wallet address", func(t *testing.T) {
		result := ge.GetActivityAcrossAllChains("")

		// Should handle empty address gracefully
		if result == "" {
			t.Error("Expected result even with empty wallet address")
		}

		t.Logf("Empty address response: %s", result)
	})

	t.Run("Test with invalid wallet address format", func(t *testing.T) {
		result := ge.GetActivityAcrossAllChains("invalid_address_12345")

		// Should handle invalid address gracefully
		if result == "" {
			t.Error("Expected result even with invalid wallet address")
		}

		t.Logf("Invalid address response: %s", result)
	})

	t.Run("Test with invalid chain name", func(t *testing.T) {
		result := ge.GetTokenBalanceForAddress("invalid_chain_12345", "0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6")

		if result == "" {
			t.Error("Expected result even with invalid chain name")
		}

		t.Logf("Invalid chain response: %s", result)
	})
}
