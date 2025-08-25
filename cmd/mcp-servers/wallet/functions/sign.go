package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"context"
	"errors"
	"fmt"
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
