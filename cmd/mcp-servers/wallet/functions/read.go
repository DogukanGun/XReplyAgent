package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (wf *WalletFunctions) ReadUserWallet() (string, error) {
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	var user []User
	ack := mg.Read(wf.MongoConnection, bson.D{{Key: "twitter_id", Value: wf.TwitterId}}, &user)
	if ack && len(user) > 0 {
		return user[0].PublicKey, nil
	}
	return "", errors.New("user not found")
}

// GenerateReadWalletTool creates an MCP tool for reading a wallet
func (wf *WalletFunctions) GenerateReadWalletTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("read_wallet",
		mcp.WithDescription("Read a wallet's public address for the current Twitter ID"),
		mcp.WithString("twitter_id", mcp.Required(), mcp.Description("Twitter id of the user")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		twitterId, _ := request.RequireString("twitter_id")
		wf.TwitterId = twitterId
		publicKey, err := wf.ReadUserWallet()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to read wallet: %v", err)), nil
		}
		return mcp.NewToolResultText(publicKey), nil
	}

	return tool, handler
}
