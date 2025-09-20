package chainlink

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"time"

	aggregator "cg-mentions-bot/cmd/mcp-servers/general/datafeeds/chainlink/aggregatorv3"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var chainRpc = map[string]string{
	"bnb":      "https://bsc.drpc.org",
	"ethereum": "https://ethereum-rpc.publicnode.com",
}

var chainIds = map[string]string{
	"bnb":      "56",
	"ethereum": "1",
}

func GetPriceFromChainlink(chainName string, assetName string) (float64, error) {
	// Read the JSON file
	var feedAddress string
	absPath, err := filepath.Abs("data_feeds.json")
	fmt.Println("Full path:", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var root Root
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		log.Fatal(err)
	}

	chainId := chainIds[chainName]

	// Search for the matching assetName and chainId
	for _, item := range root.Props.PageProps.AllFeeds {
		if item.AssetName == assetName && item.ChainID == chainId {
			feedAddress = item.ProxyAddress
			break
		}
	}

	rpcUrl := chainRpc[chainName]
	// Initialize client instance using the rpcUrl.
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Test if it is a contract address.
	ok := isContractAddress(feedAddress, client)
	if !ok {
		log.Fatalf("address %s is not a contract address\n", feedAddress)
	}

	chainlinkPriceFeedProxyAddress := common.HexToAddress(feedAddress)
	chainlinkPriceFeedProxy, err := aggregator.NewAggregatorV3Interface(chainlinkPriceFeedProxyAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	roundData, err := chainlinkPriceFeedProxy.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := chainlinkPriceFeedProxy.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	// Compute a big.int which is 10**decimals.
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Return the float value of the price
	price := divideBigInt(roundData.Answer, divisor)
	priceFloat, _ := price.Float64()
	return priceFloat, nil
}

func isContractAddress(addr string, client *ethclient.Client) bool {
	if len(addr) == 0 {
		log.Fatal("feedAddress is empty.")
	}

	// Ensure it is an Ethereum address: 0x followed by 40 hexadecimal characters.
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(addr) {
		log.Fatalf("address %s non valid\n", addr)
	}

	// Ensure it is a contract address.
	address := common.HexToAddress(addr)
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}
	isContract := len(bytecode) > 0
	return isContract
}

func formatTime(timestamp *big.Int) time.Time {
	timestampInt64 := timestamp.Int64()
	if timestampInt64 == 0 {
		log.Fatalf("timestamp %v cannot be represented as int64", timestamp)
	}
	return time.Unix(timestampInt64, 0)
}

func divideBigInt(num1 *big.Int, num2 *big.Int) *big.Float {
	if num2.BitLen() == 0 {
		log.Fatal("cannot divide by zero.")
	}
	num1BigFloat := new(big.Float).SetInt(num1)
	num2BigFloat := new(big.Float).SetInt(num2)
	result := new(big.Float).Quo(num1BigFloat, num2BigFloat)
	return result
}
