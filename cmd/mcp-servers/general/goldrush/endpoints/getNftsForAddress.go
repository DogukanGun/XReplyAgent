package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getNftsForAddress(chainName string, walletAddress string) string {

	url := ge.BaseUrl + "hyperevm-mainnet" + "/address/" + walletAddress + "/balances_nft/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateNftsForAddressTool creates an MCP tool for the getNftsForAddress endpoint
func (ge *GoldrushEndpoints) GenerateNftsForAddressTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getNftsForAddress
	nftsForAddressTool := mcp.NewTool("get_nfts_for_address",
		mcp.WithDescription("Get NFT balances for a wallet address on a specific chain"),
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
		result := ge.getNftsForAddress(chainName, walletAddress)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return nftsForAddressTool, handler
}
