package functions

import "os"

// / WalletRpc chain id to rpc
var WalletRpc = map[string]string{
	"204": os.Getenv("BNB_OP_MAINNET"),
	"97":  os.Getenv("BNB_OP_TESTNET"),
}
