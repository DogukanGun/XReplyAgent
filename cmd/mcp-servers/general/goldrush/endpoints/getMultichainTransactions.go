package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getMultichainTransactions(walletAddress string) string {

	url := ge.BaseUrl + "allchains/transactions/" + walletAddress + "/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateMultichainTransactionsTool creates an MCP tool for the getMultichainTransactions endpoint
func (ge *GoldrushEndpoints) GenerateMultichainTransactionsTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getMultichainTransactions
	multichainTransactionsTool := mcp.NewTool("get_multichain_transactions",
		mcp.WithDescription("Get transactions across all chains for a wallet address"),
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
		result := ge.getMultichainTransactions(walletAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return multichainTransactionsTool, handler
}
