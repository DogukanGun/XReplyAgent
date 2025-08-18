package endpoints

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

func (ge *GoldrushEndpoints) getTransaction(chainName string, txHash string) string {

	url := ge.BaseUrl + chainName + "/transaction_v2/" + txHash + "/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+ge.AuthToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)

}

// GenerateTransactionTool creates an MCP tool for the getTransaction endpoint
func (ge *GoldrushEndpoints) GenerateTransactionTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	// Define the tool for getTransaction
	transactionTool := mcp.NewTool("get_transaction",
		mcp.WithDescription("Get transaction details for a specific transaction hash on a specific chain"),
		mcp.WithString("chain_name",
			mcp.Required(),
			mcp.Description("Chain name to query"),
		),
		mcp.WithString("tx_hash",
			mcp.Required(),
			mcp.Description("Transaction hash to query"),
		),
	)

	// Define the handler function
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get chain name from the arguments
		chainName, err := request.RequireString("chain_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Get transaction hash from the arguments
		txHash, err := request.RequireString("tx_hash")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Call the endpoint method
		result := ge.getTransaction(chainName, txHash)

		// Return the response
		return mcp.NewToolResultText(result), nil
	}

	return transactionTool, handler
}
