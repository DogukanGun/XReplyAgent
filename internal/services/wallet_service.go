package services

import (
	"cg-mentions-bot/internal/utils/db"
	"cg-mentions-bot/internal/utils/wallet"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// WalletUser represents the user wallet data structure
type WalletUser struct {
	TwitterID        string `bson:"twitter_id" json:"twitter_id"`
	EthPublicKey     string `bson:"eth_public_key" json:"eth_public_key"`
	EthPrivateKey    string `bson:"eth_private_key" json:"eth_private_key"`
	SolanaPublicKey  string `bson:"solana_public_key" json:"solana_public_key"`
	SolanaPrivateKey string `bson:"solana_private_key" json:"solana_private_key"`
}

// WalletService handles wallet operations
type WalletService struct {
	mongoClient *mongo.Client
}

// NewWalletService creates a new wallet service
func NewWalletService(mongoClient *mongo.Client) *WalletService {
	return &WalletService{mongoClient: mongoClient}
}

// CreateOrGetWallet creates new wallets or returns existing ones for a Twitter ID
func (ws *WalletService) CreateOrGetWallet(twitterID string) (*wallet.WalletKeys, error) {
	// Check if user already has wallets
	existingWallet, err := ws.GetWallet(twitterID)
	if err == nil && existingWallet != nil {
		return existingWallet, nil
	}

	// Generate new wallets
	walletKeys, err := wallet.GenerateBothWallets()
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallets: %w", err)
	}

	// Save to database
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}

	user := WalletUser{
		TwitterID:        twitterID,
		EthPublicKey:     walletKeys.EthWallet.PublicAddress,
		EthPrivateKey:    walletKeys.EthWallet.PrivateKey,
		SolanaPublicKey:  walletKeys.SolanaWallet.PublicAddress,
		SolanaPrivateKey: walletKeys.SolanaWallet.PrivateKey,
	}

	ack := mg.Insert(ws.mongoClient, user)
	if !ack {
		return nil, fmt.Errorf("failed to save wallet to database")
	}

	return walletKeys, nil
}

// GetWallet retrieves existing wallet for a Twitter ID
func (ws *WalletService) GetWallet(twitterID string) (*wallet.WalletKeys, error) {
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}

	collection := ws.mongoClient.Database(mg.Database).Collection(mg.Collection)
	filter := bson.M{"twitter_id": twitterID}

	var user WalletUser
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found, not an error
		}
		return nil, fmt.Errorf("failed to find user wallet: %w", err)
	}

	// Check if we have the new format with both chains
	if user.EthPublicKey != "" && user.SolanaPublicKey != "" {
		return &wallet.WalletKeys{
			EthWallet: wallet.WalletKeyPair{
				PublicAddress: user.EthPublicKey,
				PrivateKey:    user.EthPrivateKey,
			},
			SolanaWallet: wallet.WalletKeyPair{
				PublicAddress: user.SolanaPublicKey,
				PrivateKey:    user.SolanaPrivateKey,
			},
		}, nil
	}

	// If we only have old format, return nil so new wallets are generated
	return nil, nil
}