package handlers

import (
	"cg-mentions-bot/internal/utils/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ContextKey string

const UidKey ContextKey = "uid"

// Database connection
var DBClient *mongo.Client
var UserCollection = db.MongoDB{
	Database:   "xreplyagent",
	Collection: "users",
}

// User model for database
type User struct {
	FirebaseID       string `bson:"firebase_id" json:"firebase_id"`
	TwitterID        string `bson:"twitter_id" json:"twitter_id"`
	Username         string `bson:"username" json:"username"`
	EthPublicKey     string `bson:"eth_public_key,omitempty" json:"eth_public_key,omitempty"`
	EthPrivateKey    string `bson:"eth_private_key,omitempty" json:"eth_private_key,omitempty"`
	SolanaPublicKey  string `bson:"solana_public_key,omitempty" json:"solana_public_key,omitempty"`
	SolanaPrivateKey string `bson:"solana_private_key,omitempty" json:"solana_private_key,omitempty"`
}

type CheckUserResponse struct {
	Register bool `json:"register"`
}

type RegisterUserRequest struct {
	TwitterID string `json:"twitter_id"`
	Username  string `json:"username"`
}

type WalletKeys struct {
	EthWallet    WalletKeyPair `json:"eth_wallet"`
	SolanaWallet WalletKeyPair `json:"solana_wallet"`
}

type WalletKeyPair struct {
	PublicAddress string `json:"public_address"`
	PrivateKey    string `json:"private_key"`
}

type RegisterUserResponse struct {
	UID       string     `json:"uid"`
	TwitterID string     `json:"twitter_id"`
	Username  string     `json:"username"`
	Message   string     `json:"message"`
	Wallets   WalletKeys `json:"wallets"`
	// Backward compatibility
	PrivateKey string `json:"private_key,omitempty"`
}

type ExecuteAppRequest struct {
	AppName string                 `json:"app_name"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type ExecuteAppResponse struct {
	UID     string `json:"uid"`
	AppName string `json:"app_name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProfileResponse struct {
	UID       string `json:"uid"`
	Username  string `json:"username,omitempty"`
	TwitterID string `json:"twitter_id,omitempty"`
	Message   string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
