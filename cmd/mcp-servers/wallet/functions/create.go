package functions

import (
	"cg-mentions-bot/internal/utils/db"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func (wf *WalletFunctions) CreateWallet() (string, error) {
	//Check if user has already session
	pk, err := wf.ReadUserWallet()
	if err == nil && pk != "" {
		return pk, nil
	}
	mg := db.MongoDB{
		Database:   "User",
		Collection: "Wallet",
	}
	//If user is not exist, create one
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate wallet: %w", err)
	}
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	user := User{
		PublicKey:  publicAddress,
		PrivateKey: privateKeyHex,
		TwitterId:  wf.TwitterId,
	}
	ack := mg.Insert(wf.MongoConnection, user)
	if !ack {
		return "", fmt.Errorf("failed to insert wallet")
	}
	return publicAddress, nil
}
