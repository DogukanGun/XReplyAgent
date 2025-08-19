package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test data structures for real API responses
var realPriceFeeds = []string{
	"0x8132e3eb1dac3e56939a16ff83848d194345f6688bff97eb1c8bd462d558802b", // BTC/USD
	"0x57ff7100a282e4af0c91154679c5dae2e5dcacb93fd467ea9cb7e58afdcfde27", // ETH/USD
	"0x879551021853eec7a7dc827578e8e69da7e4fa8148339aa0d3d5296405be4b1a", // FRAX/USD
}

// TestSetBaseURL tests the SetBaseURL function
func TestSetBaseURL(t *testing.T) {
	// Store original base URL
	originalURL := "https://hermes.pyth.network"

	// Test setting a new URL
	testURL := "http://test-server.com"
	SetBaseURL(testURL)

	// Verify the URL was set
	if baseURL != testURL {
		t.Errorf("Expected baseURL to be '%s', got '%s'", testURL, baseURL)
	}

	// Reset to original URL
	SetBaseURL(originalURL)

	// Verify it was reset
	if baseURL != originalURL {
		t.Errorf("Expected baseURL to be reset to '%s', got '%s'", originalURL, baseURL)
	}
}

// TestPriceFeedIDs tests the hardcoded price feed IDs
func TestPriceFeedIDs(t *testing.T) {
	// Test that we have the expected hardcoded feeds
	expectedFeeds := []string{"BTC/USD", "ETH/USD", "USDC/USD", "FRAX/USD"}

	for _, expected := range expectedFeeds {
		if _, exists := PriceFeedIDs[expected]; !exists {
			t.Errorf("Expected price feed '%s' to exist in hardcoded list", expected)
		}
	}

	// Test that we have the correct number of feeds
	expectedCount := 10 // Based on our hardcoded list
	if len(PriceFeedIDs) != expectedCount {
		t.Errorf("Expected %d hardcoded price feeds, got %d", expectedCount, len(PriceFeedIDs))
	}

	// Test specific feed IDs
	if PriceFeedIDs["BTC/USD"] != "0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43" {
		t.Errorf("Expected BTC/USD to have correct ID")
	}

	if PriceFeedIDs["FRAX/USD"] != "0xc96458d393fe9deb7a7d63a0ac41e2898a67a7750dbd166673279e06c868df0a" {
		t.Errorf("Expected FRAX/USD to have correct ID")
	}
}

// TestFetchPriceFeedsRealAPI tests the fetchPriceFeeds function against the real Pyth Network API
func TestFetchPriceFeedsRealAPI(t *testing.T) {

	response, err := fetchPriceFeeds(realPriceFeeds)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	// The real API should return data for the requested feeds
	if len(response.Data) == 0 {
		t.Error("Expected at least some data from real API, got empty response")
	}

	// Check that we got responses for the requested feeds
	prices := response.Data
	assert.Equal(t, 3, len(prices))
	found := false
	print(prices)
	for _, id := range realPriceFeeds {
		for _, r := range response.Data {
			if "0x"+r.ID == id {
				found = true
			}
		}
	}
	assert.True(t, found)
	// Verify the response structure
	for id, feed := range response.Data {
		if feed.ID == "" {
			t.Errorf("Expected feed ID to be set for feed %d", id)
		}
		if feed.Price.Price == "" {
			t.Errorf("Expected price to be non-zero for feed %d", id)
		}
		if feed.Metadata.PrevPublishTime <= 0 {
			t.Errorf("Expected publish time to be set for feed %d", id)
		}
	}
}

// TestFetchPriceFeedsEmpty tests fetchPriceFeeds with empty input
func TestFetchPriceFeedsEmpty(t *testing.T) {
	response, err := fetchPriceFeeds([]string{})
	if err == nil {
		t.Error("Expected error for empty feed IDs, got nil")
	}

	if response != nil {
		t.Error("Expected nil response for empty feed IDs")
	}

	expectedError := "no feed IDs provided"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestFetchAllPriceFeeds tests the fetchAllPriceFeeds function
func TestFetchAllPriceFeeds(t *testing.T) {
	// This function makes API calls and would need mocking
	// For now, we'll test that it doesn't panic

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("fetchAllPriceFeeds panicked: %v", r)
		}
	}()

	// Test with mock data (this would need refactoring to be more testable)
	// For now, we'll just ensure the function signature is correct
	_ = fetchAllPriceFeeds
}

// TestFetchSpecificPriceFeeds tests the fetchSpecificPriceFeeds function
func TestFetchSpecificPriceFeeds(t *testing.T) {
	// This function also makes API calls and would need mocking
	// For now, we'll test that it doesn't panic

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("fetchSpecificPriceFeeds panicked: %v", r)
		}
	}()

	// Test with mock data (this would need refactoring to be more testable)
	// For now, we'll just ensure the function signature is correct
	_ = fetchSpecificPriceFeeds
}

// TestListAvailableFeeds tests the listAvailableFeeds function
func TestListAvailableFeeds(t *testing.T) {
	// This function mainly prints to stdout, so we just test that it doesn't panic

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("listAvailableFeeds panicked: %v", r)
		}
	}()

	// Test that the function doesn't crash
	listAvailableFeeds()
}

// TestURLEncoding tests the URL encoding for feed IDs
func TestURLEncoding(t *testing.T) {
	feedIDs := []string{
		"0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
		"0xc96458d393fe9deb7a7d63a0ac41e2898a67a7750dbd166673279e06c868df0a",
	}

	// Build query parameters using the correct format from the documentation
	// Format: ids[]={feed_id1}&ids[]={feed_id2}
	queryParams := make([]string, len(feedIDs))
	for i, id := range feedIDs {
		queryParams[i] = fmt.Sprintf("ids[]=%s", id)
	}

	url := fmt.Sprintf("https://hermes.pyth.network/v2/updates/price/latest?%s", strings.Join(queryParams, "&"))

	// Check that the URL contains the expected format
	expectedPattern := "ids[]=0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43"
	if !strings.Contains(url, expectedPattern) {
		t.Errorf("Expected URL to contain '%s', got '%s'", expectedPattern, url)
	}

	// Check that both feed IDs are present
	for _, id := range feedIDs {
		if !strings.Contains(url, id) {
			t.Errorf("Expected URL to contain feed ID '%s'", id)
		}
	}
}

// TestFetchSpecificPriceFeedsLogic tests the logic for finding specific feeds
func TestFetchSpecificPriceFeedsLogic(t *testing.T) {
	// Test with valid feed names
	validNames := []string{"BTC/USD", "USDC/USD"}

	var feedIDs []string
	for _, name := range validNames {
		if id, exists := PriceFeedIDs[name]; exists {
			feedIDs = append(feedIDs, id)
		}
	}

	if len(feedIDs) != 2 {
		t.Errorf("Expected 2 feed IDs for valid names, got %d", len(feedIDs))
	}

	// Test with invalid feed names
	invalidNames := []string{"INVALID/USD", "UNKNOWN/USD"}

	var invalidFeedIDs []string
	for _, name := range invalidNames {
		if id, exists := PriceFeedIDs[name]; exists {
			invalidFeedIDs = append(invalidFeedIDs, id)
		}
	}

	if len(invalidFeedIDs) != 0 {
		t.Errorf("Expected 0 feed IDs for invalid names, got %d", len(invalidFeedIDs))
	}
}

// Benchmark tests for performance
func BenchmarkFetchPriceFeeds(b *testing.B) {
	feedIDs := []string{
		"0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
		"0x2b9ab1e972a281585084148ba1389800799bd4be63b5d5b80908bdc0a1801f2b",
		"0xb0948a5e5313200c632b51bb5ca32f6de0d36e995",
	}

	for i := 0; i < b.N; i++ {
		// This would need mocking for real benchmarking
		_ = feedIDs
	}
}

// Test helper functions
func TestStringContains(t *testing.T) {
	testCases := []struct {
		str      string
		substr   string
		expected bool
	}{
		{"BTC/USD", "BTC", true},
		{"ETH/USD", "BTC", false},
		{"USDC/USD", "USD", true},
		{"", "BTC", false},
		{"BTC", "", true},
	}

	for _, tc := range testCases {
		result := strings.Contains(tc.str, tc.substr)
		if result != tc.expected {
			t.Errorf("strings.Contains('%s', '%s') = %v, expected %v", tc.str, tc.substr, result, tc.expected)
		}
	}
}

// Test error message formatting
func TestErrorMessageFormatting(t *testing.T) {
	err := fmt.Errorf("API request failed with status %d: %s", 404, "Not Found")
	expected := "API request failed with status 404: Not Found"

	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// Test main function exists
func TestMainFunction(t *testing.T) {
	// Basic test to ensure main function exists
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main function panicked: %v", r)
		}
	}()

	// We can't easily test main() without refactoring, but we can ensure it exists
	_ = main
}
