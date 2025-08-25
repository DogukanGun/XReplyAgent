package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (wf *WalletFunctions) ReadUserWallet() (string, error) {
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	var user []User
	ack := mg.Read(wf.MongoConnection, bson.D{{Key: "twitter_id", Value: wf.TwitterId}}, &user)
	if ack {
		return user[0].PublicKey, nil
	}
	return "", errors.New("user not found")
}
