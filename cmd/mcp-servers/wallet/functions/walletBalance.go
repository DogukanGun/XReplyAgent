package functions

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

func (wf *WalletFunctions) GetWalletBalance() (*big.Int, error) {
	// fetch private key / address
	publicKey, err := wf.ReadUserWallet()
	if err != nil {
		return nil, fmt.Errorf("user does not exist: %w", err)
	}
	fromAddr := common.HexToAddress(publicKey)

	// connect
	ctx := context.Background()
	defer ctx.Done()
	client, err := ethclient.DialContext(ctx, os.Getenv("BNB_RPC"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// query latest balance
	balance, err := client.BalanceAt(ctx, fromAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return balance, nil
}
