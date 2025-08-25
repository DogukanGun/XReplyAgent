package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"go.mongodb.org/mongo-driver/v2/bson"
	"math/big"
	"os"
	_ "strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (wf *WalletFunctions) SignTransaction(chainId string, toAddr string, data []byte, value *big.Int) (string, error) {
	ctx := context.Background()

	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	var user []User
	ack := mg.Read(wf.MongoConnection, bson.D{{Key: "twitter_id", Value: wf.TwitterId}}, &user)
	if !ack {
		return "", errors.New("failed to find user")
	}
	privKeyHex := user[0].PrivateKey
	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Connect to RPC
	client, err := ethclient.DialContext(ctx, os.Getenv("BNB_RPC"))
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// Parse chainId string into big.Int
	chainIdInt, ok := new(big.Int).SetString(chainId, 10)
	if !ok {
		return "", fmt.Errorf("invalid chainId: %s", chainId)
	}

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Suggest EIP-1559 fees
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas tip cap: %w", err)
	}
	gasFeeCap, err := client.SuggestGasPrice(ctx) // fallback
	if err != nil {
		return "", fmt.Errorf("failed to get gas fee cap: %w", err)
	}

	// Estimate gas limit
	toAddress := common.HexToAddress(toAddr)
	msg := ethereum.CallMsg{
		From:  fromAddr,
		To:    &toAddress,
		Value: value,
		Data:  data,
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Build tx
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainIdInt,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})

	// Sign tx
	signer := types.LatestSignerForChainID(chainIdInt)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: %w", err)
	}

	// Send tx
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send tx: %w", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (wf *WalletFunctions) GenerateSignTransactionTool() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool("sign_transaction",
		mcp.WithDescription("Sign and submit a transaction"),
		mcp.WithString("chain_id", mcp.Required(), mcp.Description("Chain ID to use")),
		mcp.WithString("to_address", mcp.Required(), mcp.Description("Destination address")),
		mcp.WithString("data", mcp.Description("Hex-encoded transaction data. If you don't have any data, send empty byte array")),
		mcp.WithString("value", mcp.Description("Amount of native token to send (in wei). If you don't send anything, use 0")),
		mcp.WithString("twitter_id", mcp.Required(), mcp.Description("Twitter id of the user")),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		chainID, _ := request.RequireString("chain_id")
		toAddr, _ := request.RequireString("to_address")
		dataStr, _ := request.RequireString("data")
		valueStr, _ := request.RequireString("value")
		twitterId, _ := request.RequireString("twitter_id")
		wf.TwitterId = twitterId
		var data []byte
		if dataStr != "" {
			data = common.FromHex(dataStr)
		}
		var value *big.Int
		if valueStr != "" {
			v, ok := new(big.Int).SetString(valueStr, 10)
			if !ok {
				return mcp.NewToolResultError("invalid value parameter"), nil
			}
			value = v
		} else {
			value = big.NewInt(0)
		}

		txHash, err := wf.SignTransaction(chainID, toAddr, data, value)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(txHash), nil
	}
	return tool, handler
}
