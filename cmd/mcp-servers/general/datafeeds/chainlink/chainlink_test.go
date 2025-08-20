package chainlink

import (
	"testing"
)

func TestGetPriceFromChainlink(t *testing.T) {
	// Call the function with 'Lido Staked ETH' and 'bnb'
	price, err := GetPriceFromChainlink("ethereum", "Lido Staked ETH")

	// Check if the price is greater than zero
	if price <= 0 {
		t.Errorf("Expected price to be greater than zero, got %v", price)
	}
	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
}
