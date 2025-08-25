package tests

import (
	"cg-mentions-bot/cmd/mcp-servers/wallet/functions"
	"cg-mentions-bot/internal/utils/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestRead(t *testing.T) {
	//Arrange
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	client, err := db.ConnectToDB("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	wf := functions.WalletFunctions{
		MongoConnection: client,
		TwitterId:       RandomString(15),
	}
	pk, err := wf.CreateWallet()
	if err != nil {
		t.Failed()
	}
	var user []functions.User

	//Act
	ack := mg.Read(wf.MongoConnection, bson.D{{Key: "twitter_id", Value: wf.TwitterId}}, &user)

	//Assert
	assert.Equal(t, true, ack)
	assert.Equal(t, pk, user[0].PublicKey)

}
