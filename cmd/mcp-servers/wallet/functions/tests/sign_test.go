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
	"time"
)

func TestSignTransaction(t *testing.T) {
	//Arrange
	bnbChainTestId := "97"
	testTwitterId := "341c212cb93b395"
	toAddress := "0x5A2D55362b3ce1Bb5434c16a2aBd923c429a3446"

	mongoClient, err := db.ConnectToDB("mongodb://localhost:27017")
	defer mongoClient.Disconnect(context.Background())
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	wf := functions.WalletFunctions{
		MongoConnection: mongoClient,
		TwitterId:       testTwitterId,
	}
	ctx := context.Background()
	defer ctx.Done()
	client, err := ethclient.DialContext(ctx, os.Getenv("BNB_RPC"))

	//Act
	amount := big.NewInt(0).Mul(big.NewInt(1e14), big.NewInt(1)) // 0.0001 BNB in wei
	txHash, err := wf.SignTransaction(bnbChainTestId, toAddress, []byte{}, amount)
	if err != nil {
		t.Fatalf("failed to sign transaction: %v", err)
	}

	//Assert
	time.Sleep(8 * time.Second) // wait between polls
	receipt, _ := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	assert.Equal(t, receipt.TxHash.Hex(), txHash)
	assert.Equal(t, receipt.Status, uint64(1))
}
