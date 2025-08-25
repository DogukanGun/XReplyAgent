package tests

import (
	"cg-mentions-bot/cmd/mcp-servers/wallet/functions"
	"cg-mentions-bot/internal/utils/db"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func TestTransferAsset(t *testing.T) {
	bnbChainTestId := "97"
	testTwitterId := "341c212cb93b395"
	ctx := context.Background()
	toAddress := "0x5A2D55362b3ce1Bb5434c16a2aBd923c429a3446"
	client, err := ethclient.DialContext(ctx, os.Getenv("BNB_RPC"))
	mongoClient, err := db.ConnectToDB("mongodb://localhost:27017")
	defer mongoClient.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("failed to connect to RPC: %v", err)
	}
	wf := functions.WalletFunctions{
		MongoConnection: mongoClient,
		TwitterId:       testTwitterId,
	}
	// Act: transfer 0.0001 BNB
	amount := big.NewInt(1e14) // 0.0001 BNB
	txHash, err := wf.TransferAsset(bnbChainTestId, toAddress, amount)
	if err != nil {
		t.Fatalf("TransferAsset failed: %v", err)
	}

	// Assert
	receipt, _ := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	assert.Equal(t, uint64(1), receipt.Status, "transaction should succeed")
}
