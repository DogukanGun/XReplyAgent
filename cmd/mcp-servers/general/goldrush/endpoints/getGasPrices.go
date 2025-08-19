package endpoints

import (
	"context"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) GetGasPrices(chainName string, eventType string) string {

	url := ge.BaseUrl + "hyperevm-mainnet" + "/event/" + eventType + "/gas_prices/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body)
}

// GenerateGasPriceTool creates an MCP tool for the getGasPrices endpoint
func (ge *GoldrushEndpoints) GenerateGasPriceTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getGasPrices
	activityTool := mcp.NewTool("get_gas_prices",
		mcp.WithDescription("Commonly used to render the transfer-in and transfer-out of a token along with historical prices from an address."),
		mcp.WithString("chain_name",
			mcp.Required(),
			mcp.Description("Chain name to query"),
		),
		mcp.WithString("event_type",
			mcp.Required(),
			mcp.Description("Event type"),
			mcp.Enum("erc20", "uniswapv3", "nativetokens"),
		),
	)

	// Define the handler function
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get chain name from the arguments
		chainName, err := request.RequireString("chain_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get event type from the arguments
		eventType, err := request.RequireString("event_type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the endpoint method
		result := ge.GetGasPrices(chainName, eventType)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return activityTool, handler
}
