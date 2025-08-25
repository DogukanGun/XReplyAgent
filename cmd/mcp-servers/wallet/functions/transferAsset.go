package functions

import "math/big"

func (wf *WalletFunctions) TransferAsset(chainId string, toAddr string, amount *big.Int) (string, error) {
	// native transfer → just a tx with empty data
	return wf.SignTransaction(chainId, toAddr, []byte{}, amount)
}
