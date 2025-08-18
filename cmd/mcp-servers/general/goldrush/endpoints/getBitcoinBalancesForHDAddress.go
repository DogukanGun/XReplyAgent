package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getBitcoinBalancesForHDAddress(walletAddress string) string {

	url := ge.BaseUrl + "btc-mainnet/address/" + walletAddress + "/hd_wallets/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateBitcoinBalanceTool creates an MCP tool for the getBitcoinBalancesForHDAddress endpoint
func (ge *GoldrushEndpoints) GenerateBitcoinBalanceTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getBitcoinBalancesForHDAddress
	activityTool := mcp.NewTool("get_bitcoin_balances_for_HD_address",
		mcp.WithDescription("Fetch balances for each active child address derived from a Bitcoin HD wallet."),
		mcp.WithString("wallet_address",
			mcp.Required(),
			mcp.Description("Wallet address to query"),
		),
	)

	// Define the handler function
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get wallet address from the arguments
		walletAddress, err := request.RequireString("wallet_address")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the endpoint method
		result := ge.getBitcoinBalancesForHDAddress(walletAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return activityTool, handler
}
