package functions

import "go.mongodb.org/mongo-driver/v2/mongo"

type WalletFunctions struct {
	MongoConnection *mongo.Client
	TwitterId       string
}

type User struct {
	TwitterId        string `json:"twitter_id" bson:"twitter_id"`
	Username         string `json:"username" bson:"username"`
	EthPublicKey     string `json:"eth_public_key" bson:"eth_public_key"`
	EthPrivateKey    string `json:"eth_private_key" bson:"eth_private_key"`
	SolanaPublicKey  string `json:"solana_public_key" bson:"solana_public_key"`
	SolanaPrivateKey string `json:"solana_private_key" bson:"solana_private_key"`

	// Backward compatibility - keep old fields but mark as deprecated
	PublicKey  string `json:"public_key,omitempty" bson:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty" bson:"private_key,omitempty"`
}
