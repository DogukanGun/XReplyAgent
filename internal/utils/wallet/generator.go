package wallet

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
)

// WalletKeyPair represents a wallet's public address and private key
type WalletKeyPair struct {
	PublicAddress string `json:"public_address"`
	PrivateKey    string `json:"private_key"`
}

// WalletKeys represents both ETH and Solana wallet keys for a user
type WalletKeys struct {
	EthWallet    WalletKeyPair `json:"eth_wallet"`
	SolanaWallet WalletKeyPair `json:"solana_wallet"`
}

// GenerateEthereumWallet generates a new Ethereum wallet
func GenerateEthereumWallet() (*WalletKeyPair, error) {
	// Generate ECDSA private key for Ethereum
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ethereum private key: %w", err)
	}

	// Get public address from private key
	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Convert private key to hex string
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	return &WalletKeyPair{
		PublicAddress: publicAddress,
		PrivateKey:    privateKeyHex,
	}, nil
}

// GenerateSolanaWallet generates a new Solana wallet
func GenerateSolanaWallet() (*WalletKeyPair, error) {
	// Generate Ed25519 key pair for Solana
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Solana key pair: %w", err)
	}

	// Encode public key as base58 (Solana address format)
	publicAddress := base58.Encode(publicKey)

	// Encode private key as base58
	privateKeyB58 := base58.Encode(privateKey)

	return &WalletKeyPair{
		PublicAddress: publicAddress,
		PrivateKey:    privateKeyB58,
	}, nil
}

// GenerateBothWallets generates both Ethereum and Solana wallets
func GenerateBothWallets() (*WalletKeys, error) {
	ethWallet, err := GenerateEthereumWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ethereum wallet: %w", err)
	}

	solanaWallet, err := GenerateSolanaWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Solana wallet: %w", err)
	}

	return &WalletKeys{
		EthWallet:    *ethWallet,
		SolanaWallet: *solanaWallet,
	}, nil
}
