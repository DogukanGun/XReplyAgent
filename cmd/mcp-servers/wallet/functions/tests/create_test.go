package tests

import (
	"cg-mentions-bot/cmd/mcp-servers/wallet/functions"
	"cg-mentions-bot/internal/utils/db"
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
	"time"
)

func RandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	// hex encoding doubles the length (1 byte -> 2 chars)
	return hex.EncodeToString(bytes)[:length]
}
func TestCreate(t *testing.T) {
	client, err := db.ConnectToDB("mongodb://localhost:27017")

	t.Run("User can create a wallet", func(t *testing.T) {
		//Arrange
		if err != nil {
			t.Fatalf("failed to connect to mongodb: %v", err)
		}
		wf := functions.WalletFunctions{
			MongoConnection: client,
			TwitterId:       RandomString(15),
		}
		mg := db.MongoDB{
			Database:   "User",
			Collection: "Wallet",
		}
		var user []functions.User

		//Act
		pk, err := wf.CreateWallet()
		if err != nil {
			t.Failed()
		}
		time.Sleep(2 * time.Second)
		ack := mg.Read(wf.MongoConnection, bson.D{{Key: "twitter_id", Value: wf.TwitterId}}, &user)
		if ack {
			t.Failed()
		}

		//Assert
		assert.Equal(t, pk, user[0].PublicKey)

	})

	t.Cleanup(func() {
		_ = client.Disconnect(context.Background())
	})
}
