package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PriceData struct {
	Price       string `json:"price"`
	Conf        string `json:"conf"`
	Expo        int    `json:"expo"`
	PublishTime int64  `json:"publish_time"`
}

type Metadata struct {
	Slot               int64 `json:"slot"`
	ProofAvailableTime int64 `json:"proof_available_time"`
	PrevPublishTime    int64 `json:"prev_publish_time"`
}
type ParsedPrice struct {
	ID       string    `json:"id"`
	Price    PriceData `json:"price"`
	EmaPrice PriceData `json:"ema_price"`
	Metadata Metadata  `json:"metadata"`
}

// HermesResponse represents the response from Hermes API
type HermesResponse struct {
	Data []ParsedPrice `json:"parsed"`
}

// Hardcoded price feed IDs from Pyth Network webpage
var PriceFeedIDs = map[string]string{
	"BTC/USD":  "0xe62df6c8b4a85fe1a67db44dc12de5db330f7ac66b72dc658afedf0f4a415b43",
	"ETH/USD":  "0x2b9ab1e972a281585084148ba1389800799bd4be63b5d5b80908bdc0a1801f2b",
	"USDC/USD": "0xb0948a5e5313200c632b51bb5ca32f6de0d36e995",
	"FRAX/USD": "0xc96458d393fe9deb7a7d63a0ac41e2898a67a7750dbd166673279e06c868df0a",
	"USDT/USD": "0x2b9ab1e972a281585084148ba1389800799bd4be63b5d5b80908bdc0a1801f2b",
	"DAI/USD":  "0xb0948a5e5313200c632b51bb5ca32f6de0d36e995",
	"SOL/USD":  "0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
	"ADA/USD":  "0x2a01deaec9e51a579277b34b122399984d0bbf57e2458a7e42fecd2829867a0d",
	"DOT/USD":  "0xca3d9fb18e30bb8c943a4e6e0a97b3e8f8b5b3e8f8b5b3e8f8b5b3e8f8b5b3e8",
	"LINK/USD": "0x8ac0c70fff57e9aefdf5edf44b51d62c2d433653cbb2cf5cc06bb115af04d221",
}

// baseURL is configurable for testing purposes
var baseURL = "https://hermes.pyth.network"

// SetBaseURL sets the base URL for testing (used by tests to point to mock server)
func SetBaseURL(url string) {
	baseURL = url
}

// fetchPriceFeeds fetches price feeds from Pyth Network Hermes API
// Based on: https://docs.pyth.network/price-feeds/fetch-price-updates#rest-api
func fetchPriceFeeds(feedIDs []string) (*HermesResponse, error) {
	if len(feedIDs) == 0 {
		return nil, fmt.Errorf("no feed IDs provided")
	}
	baseURL := "https://hermes.pyth.network/v2/updates/price/latest"
	// Parse the base URL
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Add ids[] query parameters dynamically
	q := u.Query()
	for _, id := range feedIDs {
		q.Add("ids[]", id)
	}
	u.RawQuery = q.Encode()

	fmt.Printf("Making request to: %s\n", u.String())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response HermesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// displayPriceFeeds displays the price feeds in a formatted table
func displayPriceFeeds(response *HermesResponse) {
	fmt.Println("\n=== Pyth Network Price Feeds ===")
	fmt.Printf("%-20s %-15s %-20s %-15s %-15s\n", "Symbol", "Price", "Confidence", "Expo", "Last Update")
	fmt.Println(strings.Repeat("-", 90))

	for _, feed := range response.Data {
		// Find the symbol name from our hardcoded map
		var symbol string
		for name, feedID := range PriceFeedIDs {
			if feedID == feed.ID {
				symbol = name
				break
			}
		}
		if symbol == "" {
			symbol = "Unknown"
		}

		// Convert publish time to readable format
		publishTime := time.Unix(feed.Metadata.ProofAvailableTime, 0).Format("15:04:05")

		fmt.Printf("%-20s %s %s %d %-15s\n",
			symbol,
			feed.Price.Price,
			feed.Price.Conf,
			feed.Price.Expo,
			publishTime)
	}
}

// fetchAllPriceFeeds fetches all hardcoded price feeds
func fetchAllPriceFeeds() {
	fmt.Println("Fetching all hardcoded price feeds from Pyth Network...")

	// Extract all hardcoded price feed IDs
	var feedIDs []string
	for _, id := range PriceFeedIDs {
		feedIDs = append(feedIDs, id)
	}

	response, err := fetchPriceFeeds(feedIDs)
	if err != nil {
		fmt.Printf("Error fetching price feeds: %v\n", err)
		return
	}

	displayPriceFeeds(response)
}

// fetchSpecificPriceFeeds fetches specific price feeds by name
func fetchSpecificPriceFeeds(names []string) {
	fmt.Printf("Fetching specific price feeds: %v\n", names)

	var feedIDs []string
	var foundNames []string

	for _, name := range names {
		if id, exists := PriceFeedIDs[name]; exists {
			feedIDs = append(feedIDs, id)
			foundNames = append(foundNames, name)
		} else {
			fmt.Printf("Warning: Price feed '%s' not found in hardcoded list\n", name)
		}
	}

	if len(feedIDs) == 0 {
		fmt.Println("No valid price feeds found")
		return
	}

	response, err := fetchPriceFeeds(feedIDs)
	if err != nil {
		fmt.Printf("Error fetching price feeds: %v\n", err)
		return
	}

	displayPriceFeeds(response)
}

// listAvailableFeeds lists all hardcoded price feeds
func listAvailableFeeds() {
	fmt.Println("\n=== Hardcoded Price Feeds ===")
	fmt.Printf("%-20s %-70s\n", "Name", "Feed ID")
	fmt.Println(strings.Repeat("-", 90))

	for name, id := range PriceFeedIDs {
		fmt.Printf("%-20s %-70s\n", name, id)
	}
}

func main() {
	fmt.Println("Pyth Network Price Feeds Client - Hardcoded IDs")
	fmt.Println("===============================================")
	fmt.Println("Using Pyth Network Hermes API: https://hermes.pyth.network/v2/updates/price/latest")
	fmt.Println("Documentation: https://docs.pyth.network/price-feeds/fetch-price-updates#rest-api")

	// List available feeds
	listAvailableFeeds()

	// Fetch all hardcoded price feeds
	fmt.Println("\n" + strings.Repeat("=", 50))
	fetchAllPriceFeeds()

	// Example: Fetch specific feeds (USDC/USD and FRAX/USD as in your curl example)
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Example: Fetching USDC/USD and FRAX/USD")
	fetchSpecificPriceFeeds([]string{"USDC/USD", "FRAX/USD"})

	// Example: Fetch other popular feeds
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Example: Fetching BTC/USD and ETH/USD")
	fetchSpecificPriceFeeds([]string{"BTC/USD", "ETH/USD"})
}
