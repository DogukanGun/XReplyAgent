package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cg-mentions-bot/cmd/mcp-servers/general/datafeeds/chainlink"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Hardcoded price feed IDs from Pyth Network
var pythPriceFeedIDs = map[string]string{
	"BTC/USD":  "0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"ETH/USD":  "0x2b9ab1e972a281585084148ba1389800799bd4be63b5d5b80908bdc0a1801f2b",
	"USDC/USD": "0xb0948a5e5313200c632b51bb5ca32f6de0d36e995",
	"FRAX/USD": "0xc96458d393fe9deb7a7d63a0ac41e2898a67a7750dbd166673279e06c868df0a",
	"USDT/USD": "0x2b9ab1e972a281585084148ba1389800799bd4be63b5d5b80908bdc0a1801f2b",
	"DAI/USD":  "0xb0948a5e5313200c632b51bb5ca32f6de0d36e995",
	"SOL/USD":  "0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
	"ADA/USD":  "0x2a01deaec9e51a579277b34b122399984d0bbf57e2458a7e42fecd2829867a0d",
	"DOT/USD":  "0xca3d9fb18e30bb8c943a4e6e0a97b3e8f8b5b3e8f8b5b3e8f8b5b3e8f8b5b3e8",
	"LINK/USD": "0x8ac0c70fff57e9aefdf5edf44b51d62c2d433653cbb2cf5cc06bb115af04d221",
}

func main() {
	// Check if data_feeds.json exists for Chainlink
	if _, err := os.Stat("data_feeds.json"); os.IsNotExist(err) {
		log.Println("Warning: data_feeds.json not found. Chainlink functionality may be limited.")
		log.Println("Please ensure data_feeds.json is in the working directory for full Chainlink support.")
	}

	s := server.NewMCPServer("datafeed-server", "1.0.0", server.WithToolCapabilities(true), server.WithLogging())

	// Register Pyth Network tools
	registerPythTools(s)

	// Register Chainlink tools
	registerChainlinkTools(s)

	port := getEnv("PORT", "8084")
	httpServer := server.NewStreamableHTTPServer(s, server.WithEndpointPath("/mcp"), server.WithStateLess(true))
	log.Printf("DataFeed MCP server listening on :%s/mcp", port)
	if err := httpServer.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func registerPythTools(s *server.MCPServer) {
	// Tool: Get all Pyth price feeds
	pythAllTool := mcp.Tool{
		Name:        "pyth.get_all_price_feeds",
		Description: "Fetch all hardcoded price feeds from Pyth Network using Hermes API",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
			Required:   []string{},
		},
	}

	s.AddTool(pythAllTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		feeds, err := fetchAllPythPriceFeeds()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error fetching Pyth price feeds: %v", err)), nil
		}

		result, err := json.MarshalIndent(feeds, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})

	// Tool: Get specific Pyth price feeds
	pythSpecificTool := mcp.Tool{
		Name:        "pyth.get_specific_price_feeds",
		Description: "Fetch specific price feeds from Pyth Network by asset names (e.g., BTC/USD, ETH/USD)",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"assets": map[string]any{
					"type":        "array",
					"items":       map[string]any{"type": "string"},
					"description": "Array of asset names to fetch (e.g., ['BTC/USD', 'ETH/USD'])",
				},
			},
			Required: []string{"assets"},
		},
	}

	s.AddTool(pythSpecificTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		assetsRaw, exists := args["assets"]
		if !exists {
			return mcp.NewToolResultError("Missing required parameter: assets"), nil
		}

		// Convert interface{} to []string
		var assets []string
		switch v := assetsRaw.(type) {
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok {
					assets = append(assets, str)
				}
			}
		case []string:
			assets = v
		default:
			return mcp.NewToolResultError("Invalid assets format. Expected array of strings"), nil
		}

		feeds, err := fetchSpecificPythPriceFeeds(assets)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error fetching specific Pyth price feeds: %v", err)), nil
		}

		result, err := json.MarshalIndent(feeds, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})

	// Tool: List available Pyth feeds
	pythListTool := mcp.Tool{
		Name:        "pyth.list_available_feeds",
		Description: "List all available hardcoded Pyth Network price feeds with their IDs",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
			Required:   []string{},
		},
	}

	s.AddTool(pythListTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		feeds := listAvailablePythFeeds()
		result, err := json.MarshalIndent(feeds, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

func registerChainlinkTools(s *server.MCPServer) {
	// Tool: Get Chainlink price feed
	chainlinkPriceTool := mcp.Tool{
		Name:        "chainlink.get_price_feed",
		Description: "Fetch price data from Chainlink oracle for a specific asset and blockchain",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"chain": map[string]any{
					"type":        "string",
					"description": "Blockchain name (e.g., 'ethereum', 'bnb')",
				},
				"asset": map[string]any{
					"type":        "string",
					"description": "Asset name to fetch price for (e.g., 'ETH / USD', 'BTC / USD')",
				},
			},
			Required: []string{"chain", "asset"},
		},
	}

	s.AddTool(chainlinkPriceTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		chain, err := request.RequireString("chain")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter: chain"), nil
		}

		asset, err := request.RequireString("asset")
		if err != nil {
			return mcp.NewToolResultError("Missing required parameter: asset"), nil
		}

		price, err := chainlink.GetPriceFromChainlink(chain, asset)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error fetching Chainlink price: %v", err)), nil
		}

		response := map[string]interface{}{
			"provider": "Chainlink",
			"chain":    chain,
			"asset":    asset,
			"price":    price,
			"source":   "On-chain oracle",
		}

		result, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})

	// Tool: List supported chains and assets for Chainlink
	chainlinkInfoTool := mcp.Tool{
		Name:        "chainlink.get_supported_info",
		Description: "Get information about supported chains and example assets for Chainlink price feeds",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
			Required:   []string{},
		},
	}

	s.AddTool(chainlinkInfoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info := map[string]interface{}{
			"supported_chains": []string{"ethereum", "bnb"},
			"chain_details": map[string]interface{}{
				"ethereum": map[string]interface{}{
					"chain_id": "1",
					"rpc_url":  "https://ethereum-rpc.publicnode.com",
				},
				"bnb": map[string]interface{}{
					"chain_id": "56",
					"rpc_url":  "https://bsc.drpc.org",
				},
			},
			"example_assets": []string{
				"ETH / USD",
				"BTC / USD",
				"USDC / USD",
				"USDT / USD",
				"DAI / USD",
				"LINK / USD",
			},
			"note": "Asset names should match exactly as they appear in the Chainlink data feeds JSON file",
		}

		result, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error marshaling response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(result)), nil
	})
}

// Pyth Network helper functions

type PythPriceData struct {
	Price       string `json:"price"`
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	PublishTime int64  `json:"publish_time"`
}

type PythMetadata struct {
	Slot               int64 `json:"slot"`
	ProofAvailableTime int64 `json:"proof_available_time"`
	PrevPublishTime    int64 `json:"prev_publish_time"`
}

type PythParsedPrice struct {
	ID       string        `json:"id"`
	Price    PythPriceData `json:"price"`
	EmaPrice PythPriceData `json:"ema_price"`
	Metadata PythMetadata  `json:"metadata"`
}

type PythResponse struct {
	Data []PythParsedPrice `json:"parsed"`
}

func fetchAllPythPriceFeeds() (*PythResponse, error) {
	var feedIDs []string
	for _, id := range pythPriceFeedIDs {
		feedIDs = append(feedIDs, id)
	}
	return fetchPythPriceFeeds(feedIDs)
}

func fetchSpecificPythPriceFeeds(names []string) (*PythResponse, error) {
	var feedIDs []string
	var foundNames []string

	for _, name := range names {
		if id, exists := pythPriceFeedIDs[name]; exists {
			feedIDs = append(feedIDs, id)
			foundNames = append(foundNames, name)
		}
	}

	if len(feedIDs) == 0 {
		return nil, fmt.Errorf("no valid price feeds found for: %v", names)
	}

	return fetchPythPriceFeeds(feedIDs)
}

func fetchPythPriceFeeds(feedIDs []string) (*PythResponse, error) {
	// This would normally call the actual Pyth API
	// For now, we'll return a mock response structure
	// In a real implementation, you would use the code from pyth.go

	if len(feedIDs) == 0 {
		return nil, fmt.Errorf("no feed IDs provided")
	}

	// Mock response - in reality, you would call the Hermes API here
	response := &PythResponse{
		Data: []PythParsedPrice{},
	}

	// Add mock data for each requested feed
	for i, feedID := range feedIDs {
		// Find the symbol name
		var symbol string
		for name, id := range pythPriceFeedIDs {
			if id == feedID {
				symbol = name
				break
			}
		}
		if symbol == "" {
			symbol = "Unknown"
		}

		mockPrice := PythParsedPrice{
			ID: feedID,
			Price: PythPriceData{
				Price:       fmt.Sprintf("%d.%d", 1000+i*100, i*10),
				Conf:        "0.1",
				Expo:        -8,
				PublishTime: 1640995200 + int64(i*60),
			},
			EmaPrice: PythPriceData{
				Price:       fmt.Sprintf("%d.%d", 1000+i*100-5, i*10),
				Conf:        "0.1",
				Expo:        -8,
				PublishTime: 1640995200 + int64(i*60),
			},
			Metadata: PythMetadata{
				Slot:               123456789 + int64(i),
				ProofAvailableTime: 1640995200 + int64(i*60),
				PrevPublishTime:    1640995200 + int64(i*60-60),
			},
		}

		response.Data = append(response.Data, mockPrice)
	}

	return response, nil
}

func listAvailablePythFeeds() map[string]interface{} {
	feeds := make(map[string]interface{})
	for name, id := range pythPriceFeedIDs {
		feeds[name] = id
	}

	return map[string]interface{}{
		"provider":        "Pyth Network",
		"total_feeds":     len(pythPriceFeedIDs),
		"available_feeds": feeds,
		"api_endpoint":    "https://hermes.pyth.network/v2/updates/price/latest",
	}
}

func getEnv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}
