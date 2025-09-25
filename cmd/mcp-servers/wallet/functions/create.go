package functions

import (
	"cg-mentions-bot/internal/services"
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func (wf *WalletFunctions) CreateWallet() (string, error) {
	//Check if user has already session
	pk, err := wf.ReadUserWallet()
	if err == nil && pk != "" {
		return pk, nil
	}

	// Use the common wallet service
	walletService := services.NewWalletService(wf.MongoConnection)
	walletKeys, err := walletService.CreateOrGetWallet(wf.TwitterId)
	if err != nil {
		return "", fmt.Errorf("failed to create/get wallets: %w", err)
	}

	return walletKeys.EthWallet.PublicAddress, nil
}

func (wf *WalletFunctions) GenerateCreateWalletTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("create_wallet",
		mcp.WithDescription("Create a new wallet for a given Twitter ID, or return an existing wallet if one already exists."),
		mcp.WithString("twitter_id", mcp.Required(), mcp.Description("Twitter id of the user")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		twitterId, _ := request.RequireString("twitter_id")
		wf.TwitterId = twitterId
		publicKey, err := wf.CreateWallet()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to create wallet: %v", err)), nil
		}
		return mcp.NewToolResultText(publicKey), nil
	}

	return tool, handler
}
