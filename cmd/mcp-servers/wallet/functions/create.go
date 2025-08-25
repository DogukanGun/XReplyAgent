package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mark3labs/mcp-go/mcp"
)

func (wf *WalletFunctions) CreateWallet() (string, error) {
	//Check if user has already session
	pk, err := wf.ReadUserWallet()
	if err == nil && pk != "" {
		return pk, nil
	}
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	//If user is not exist, create one
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate wallet: %w", err)
	}
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	user := User{
		PublicKey:  publicAddress,
		PrivateKey: privateKeyHex,
		TwitterId:  wf.TwitterId,
	}
	ack := mg.Insert(wf.MongoConnection, user)
	if !ack {
		return "", fmt.Errorf("failed to insert wallet")
	}
	return publicAddress, nil
}

func (wf *WalletFunctions) GenerateCreateWalletTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("create_wallet",
		mcp.WithDescription("Create a new wallet for a given Twitter ID, or return an existing wallet if one already exists"),
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
