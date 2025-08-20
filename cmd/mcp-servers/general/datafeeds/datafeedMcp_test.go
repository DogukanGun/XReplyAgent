package datafeeds

import (
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewDataFeedServer(t *testing.T) {
	server := NewDataFeedServer()
	if server == nil {
		t.Error("Expected non-nil server")
	}

	if server.mcpServer == nil {
		t.Error("Expected non-nil MCP server")
	}
}

func TestPythPriceFeedIDs(t *testing.T) {
	expectedFeeds := []string{"BTC/USD", "ETH/USD", "USDC/USD", "FRAX/USD", "USDT/USD", "DAI/USD", "SOL/USD", "ADA/USD", "DOT/USD", "LINK/USD"}

	for _, feed := range expectedFeeds {
		if _, exists := pythPriceFeedIDs[feed]; !exists {
			t.Errorf("Expected feed %s to exist in pythPriceFeedIDs", feed)
		}
	}

	if len(pythPriceFeedIDs) != len(expectedFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(expectedFeeds), len(pythPriceFeedIDs))
	}
}

func TestFetchAllPythPriceFeeds(t *testing.T) {
	dfs := NewDataFeedServer()

	response, err := dfs.fetchAllPythPriceFeeds()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
		return
	}

	if len(response.Data) == 0 {
		t.Error("Expected at least one price feed in response")
	}

	// Check that we have the expected number of feeds
	if len(response.Data) != len(pythPriceFeedIDs) {
		t.Errorf("Expected %d feeds, got %d", len(pythPriceFeedIDs), len(response.Data))
	}
}

func TestFetchSpecificPythPriceFeeds(t *testing.T) {
	dfs := NewDataFeedServer()

	// Test with valid feeds
	validFeeds := []string{"BTC/USD", "ETH/USD"}
	response, err := dfs.fetchSpecificPythPriceFeeds(validFeeds)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
		return
	}

	if len(response.Data) != len(validFeeds) {
		t.Errorf("Expected %d feeds, got %d", len(validFeeds), len(response.Data))
	}

	// Test with invalid feeds
	invalidFeeds := []string{"INVALID/USD", "UNKNOWN/USD"}
	_, err = dfs.fetchSpecificPythPriceFeeds(invalidFeeds)
	if err == nil {
		t.Error("Expected error for invalid feeds, got nil")
	}

	// Test with mixed valid and invalid feeds
	mixedFeeds := []string{"BTC/USD", "INVALID/USD"}
	response, err = dfs.fetchSpecificPythPriceFeeds(mixedFeeds)
	if err != nil {
		t.Errorf("Expected no error for mixed feeds, got: %v", err)
	}

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 valid feed, got %d", len(response.Data))
	}
}

func TestListAvailablePythFeeds(t *testing.T) {
	dfs := NewDataFeedServer()

	feeds := dfs.listAvailablePythFeeds()
	if feeds == nil {
		t.Error("Expected non-nil feeds map")
		return
	}

	// Check required fields
	provider, exists := feeds["provider"]
	if !exists || provider != "Pyth Network" {
		t.Error("Expected provider to be 'Pyth Network'")
	}

	totalFeeds, exists := feeds["total_feeds"]
	if !exists {
		t.Error("Expected total_feeds field")
	}

	if totalFeeds != len(pythPriceFeedIDs) {
		t.Errorf("Expected total_feeds to be %d, got %v", len(pythPriceFeedIDs), totalFeeds)
	}

	availableFeeds, exists := feeds["available_feeds"]
	if !exists {
		t.Error("Expected available_feeds field")
	}

	feedsMap, ok := availableFeeds.(map[string]interface{})
	if !ok {
		t.Error("Expected available_feeds to be a map")
		return
	}

	if len(feedsMap) != len(pythPriceFeedIDs) {
		t.Errorf("Expected %d available feeds, got %d", len(pythPriceFeedIDs), len(feedsMap))
	}
}

func TestFetchPythPriceFeedsWithEmptyIDs(t *testing.T) {
	dfs := NewDataFeedServer()

	_, err := dfs.fetchPythPriceFeeds([]string{})
	if err == nil {
		t.Error("Expected error for empty feed IDs, got nil")
	}

	if err.Error() != "no feed IDs provided" {
		t.Errorf("Expected 'no feed IDs provided' error, got: %v", err)
	}
}

func TestGetEnv(t *testing.T) {
	// Test with default value
	result := getEnv("NON_EXISTENT_ENV_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", result)
	}

	// Test with existing environment variable
	t.Setenv("TEST_ENV_VAR", "test_value")
	result = getEnv("TEST_ENV_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestPythResponseStructure(t *testing.T) {
	// Test JSON marshaling/unmarshaling of Pyth structures
	mockResponse := PythResponse{
		Data: []PythParsedPrice{
			{
				ID: "test_id",
				Price: PythPriceData{
					Price:       "1000.50",
					Conf:        "0.1",
					Expo:        -8,
					PublishTime: 1640995200,
				},
				EmaPrice: PythPriceData{
					Price:       "1000.45",
					Conf:        "0.1",
					Expo:        -8,
					PublishTime: 1640995200,
				},
				Metadata: PythMetadata{
					Slot:               123456789,
					ProofAvailableTime: 1640995200,
					PrevPublishTime:    1640995140,
				},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(mockResponse)
	if err != nil {
		t.Errorf("Failed to marshal PythResponse: %v", err)
	}

	// Unmarshal back
	var unmarshaled PythResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal PythResponse: %v", err)
	}

	// Verify data integrity
	if len(unmarshaled.Data) != 1 {
		t.Errorf("Expected 1 price feed, got %d", len(unmarshaled.Data))
	}

	if unmarshaled.Data[0].ID != "test_id" {
		t.Errorf("Expected ID 'test_id', got '%s'", unmarshaled.Data[0].ID)
	}

	if unmarshaled.Data[0].Price.Price != "1000.50" {
		t.Errorf("Expected price '1000.50', got '%s'", unmarshaled.Data[0].Price.Price)
	}
}

// Mock test for MCP tool integration (would require more complex setup for full integration testing)
func TestMCPToolStructure(t *testing.T) {
	// Test that we can create a proper MCP tool request structure
	args := map[string]interface{}{
		"assets": []string{"BTC/USD", "ETH/USD"},
	}

	request := mcp.CallToolRequest{
		Request: mcp.Request{Method: "tools/call"},
		Params: mcp.CallToolParams{
			Name:      "pyth.get_specific_price_feeds",
			Arguments: args,
		},
	}

	// Test argument access
	arguments := request.GetArguments()
	if arguments == nil {
		t.Error("Expected non-nil arguments")
	}

	assets, exists := arguments["assets"]
	if !exists {
		t.Error("Expected 'assets' argument to exist")
	}

	assetsList, ok := assets.([]string)
	if !ok {
		t.Error("Expected assets to be []string")
		return
	}

	if len(assetsList) != 2 {
		t.Errorf("Expected 2 assets, got %d", len(assetsList))
	}

	if assetsList[0] != "BTC/USD" || assetsList[1] != "ETH/USD" {
		t.Errorf("Expected ['BTC/USD', 'ETH/USD'], got %v", assetsList)
	}
}

// Benchmark tests
func BenchmarkFetchAllPythPriceFeeds(b *testing.B) {
	dfs := NewDataFeedServer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := dfs.fetchAllPythPriceFeeds()
		if err != nil {
			b.Errorf("Error in benchmark: %v", err)
		}
	}
}

func BenchmarkListAvailablePythFeeds(b *testing.B) {
	dfs := NewDataFeedServer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		feeds := dfs.listAvailablePythFeeds()
		if feeds == nil {
			b.Error("Expected non-nil feeds")
		}
	}
}
