package functions

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mark3labs/mcp-go/mcp"
)

func (wf *WalletFunctions) GetWalletBalance() (*big.Int, error) {
	// fetch private key / address
	publicKey, err := wf.ReadUserWallet()
	if err != nil {
		return nil, fmt.Errorf("user does not exist: %w", err)
	}
	fromAddr := common.HexToAddress(publicKey)

	// connect
	ctx := context.Background()
	defer ctx.Done()
	client, err := ethclient.DialContext(ctx, os.Getenv("BNB_RPC"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// query latest balance
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}

func (wf *WalletFunctions) GenerateGetWalletBalanceTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("get_wallet_balance",
		mcp.WithDescription("Get the balance of the wallet"),
		mcp.WithString("twitter_id", mcp.Required(), mcp.Description("Twitter id of the user")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		twitterId, _ := request.RequireString("twitter_id")
		wf.TwitterId = twitterId

		balance, err := wf.GetWalletBalance()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(balance.String()), nil
	}
	return tool, handler
}
