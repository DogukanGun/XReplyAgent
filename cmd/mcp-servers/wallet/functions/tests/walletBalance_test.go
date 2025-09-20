package tests

import (
	"cg-mentions-bot/cmd/mcp-servers/wallet/functions"
	"cg-mentions-bot/internal/utils/db"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWalletBalance(t *testing.T) {
	//Arrange
	testTwitterId := "341c212cb93b395"
	mongoClient, err := db.ConnectToDB("mongodb://localhost:27017")
	defer mongoClient.Disconnect(context.Background())
	wf := functions.WalletFunctions{
		MongoConnection: mongoClient,
		TwitterId:       testTwitterId,
	}

	// Act
	balance, err := wf.GetWalletBalance()
	if err != nil {
		t.Fatalf("GetWalletBalance failed: %v", err)
	}

	// Assert
	assert.NotNil(t, balance)
	assert.True(t, balance.Sign() >= 0, "balance must be >= 0")
	t.Logf("Wallet balance: %s wei", balance.String())
}
