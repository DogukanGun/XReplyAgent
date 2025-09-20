package functions

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"math/big"
)

func (wf *WalletFunctions) TransferAsset(chainId string, toAddr string, amount *big.Int) (string, error) {
	// native transfer → just a tx with empty data
	return wf.SignTransaction(chainId, toAddr, []byte{}, amount)
}

func (wf *WalletFunctions) GenerateTransferAssetTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("transfer_asset",
		mcp.WithDescription("Transfer native asset to an address"),
		mcp.WithString("chain_id", mcp.Required(), mcp.Description("Chain ID to use")),
		mcp.WithString("to_address", mcp.Required(), mcp.Description("Recipient address")),
		mcp.WithString("amount", mcp.Required(), mcp.Description("Amount in wei")),
		mcp.WithString("twitter_id", mcp.Required(), mcp.Description("Twitter id of the user")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		chainID, _ := request.RequireString("chain_id")
		toAddr, _ := request.RequireString("to_address")
		amountStr, _ := request.RequireString("amount")
		twitterId, _ := request.RequireString("twitter_id")
		wf.TwitterId = twitterId

		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return mcp.NewToolResultError("invalid amount parameter"), nil
		}

		txHash, err := wf.TransferAsset(chainID, toAddr, amount)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(txHash), nil
	}
	return tool, handler
}
