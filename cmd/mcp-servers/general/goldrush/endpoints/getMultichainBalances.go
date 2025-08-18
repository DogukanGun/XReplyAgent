package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) GetMultichainBalances(walletAddress string) string {

	url := ge.BaseUrl + "allchains/address/" + walletAddress + "/balances/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateMultichainBalancesTool creates an MCP tool for the getMultichainBalances endpoint
func (ge *GoldrushEndpoints) GenerateMultichainBalancesTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getMultichainBalances
	multichainBalancesTool := mcp.NewTool("get_multichain_balances",
		mcp.WithDescription("Get balances across all chains for a wallet address"),
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
		result := ge.GetMultichainBalances(walletAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return multichainBalancesTool, handler
}
