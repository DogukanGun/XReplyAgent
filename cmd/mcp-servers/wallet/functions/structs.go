package functions

import "go.mongodb.org/mongo-driver/v2/mongo"

type WalletFunctions struct {
	MongoConnection *mongo.Client
	TwitterId       string
}

type User struct {
	TwitterId  string `json:"twitter_id" bson:"twitter_id"`
	Username   string `json:"username" bson:"username"`
	PublicKey  string `json:"public_key" bson:"public_key"`
	PrivateKey string `json:"private_key" bson:"private_key"`
}
