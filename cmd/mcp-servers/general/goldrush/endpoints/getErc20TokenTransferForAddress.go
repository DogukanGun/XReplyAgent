package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getErc20TokenTransferForAddress(chainName string, walletAddress string) string {

	url := ge.BaseUrl + chainName + "/address/" + walletAddress + "/transfers_v2/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateErc20TokenTransferTool creates an MCP tool for the getErc20TokenTransferForAddress endpoint
func (ge *GoldrushEndpoints) GenerateErc20TokenTransferTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getErc20TokenTransferForAddress
	activityTool := mcp.NewTool("get_erc20_token_transfer_for_address",
		mcp.WithDescription("Commonly used to render the transfer-in and transfer-out of a token along with historical prices from an address."),
		mcp.WithString("chain_name",
			mcp.Required(),
			mcp.Description("Chain name to query"),
		),
		mcp.WithString("wallet_address",
			mcp.Required(),
			mcp.Description("Wallet address to query"),
		),
	)

	// Define the handler function
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get chain name from the arguments
		chainName, err := request.RequireString("chain_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get wallet address from the arguments
		walletAddress, err := request.RequireString("wallet_address")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the endpoint method
		result := ge.getErc20TokenTransferForAddress(chainName, walletAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return activityTool, handler
}
