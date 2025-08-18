package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getHistoricalTokenPrices(chainName string, contractAddress string) string {

	url := ge.BaseUrl + "pricing/historical_by_addresses_v2/" + chainName + "/USD/" + contractAddress + "/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateHistoricalTokenPriceTool creates an MCP tool for the getGasPrices endpoint
func (ge *GoldrushEndpoints) GenerateHistoricalTokenPriceTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getHistoricalTokenPrices
	activityTool := mcp.NewTool("get_historical_token_prices",
		mcp.WithDescription("Commonly used to fetch the historical native and fungible (ERC20) tokens held "+
			"by an address at a given block height or date. Response includes daily prices and other metadata."),
		mcp.WithString("chain_name",
			mcp.Required(),
			mcp.Description("Chain name to query"),
		),
		mcp.WithString("contract_address",
			mcp.Required(),
			mcp.Description("The requested address. Passing in an ENS, RNS, Lens Handle, or an Unstoppable Domain resolves automatically."),
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
		contractAddress, err := request.RequireString("contract_address")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the endpoint method
		result := ge.getHistoricalTokenPrices(chainName, contractAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return activityTool, handler
}
