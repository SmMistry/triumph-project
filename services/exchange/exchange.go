package exchange

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"strings"
)

// Exchange defines an interface for interacting with cryptocurrency exchanges
type Exchange interface {
	// GetPrices retrives the buy and sell prices from an exchange
	// It takes a context and symbol returning a buy price, sell price, error 
	GetPrices(ctx context.Context, symbol string) (float64, float64, error)
	// Get the name of the current exchange
	GetName() string
}

// CoinbaseExchange implements the Exchange interface for Coinbase
type CoinbaseExchange struct{}

// GetPrices retrieves the price for a given symbol from Coinbase
func (c *CoinbaseExchange) GetPrices(ctx context.Context, symbol string) (float64, float64, error) {
	// Construct the Coinbase API URL
	// url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s/spot", symbol)
	url := fmt.Sprintf("https://api.exchange.coinbase.com/products/%s-USD/book", symbol)

	// Create a new HTTP client with a timeout
	client := http.Client{Timeout: 10 * time.Second}

	// Send the request to the Coinbase API
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get price from coinbase: %w", err)
	}
	defer resp.Body.Close()

	// Decode the JSON response
	// The only fields we're intrested in are the two prices
	var coinbaseResponse struct {
		Bids [][]any `json:"bids"`
		Asks [][]any `json:"asks"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&coinbaseResponse); err != nil {
		return 0, 0, fmt.Errorf("failed to decode coinbase response: %w", err)
	}

	// Make sure bid and ask data are present
	if len(coinbaseResponse.Bids) == 0 || len(coinbaseResponse.Bids[0]) == 0 {
		return 0, 0, fmt.Errorf("Failed to find bid prices in coinbase response")
	}
	if len(coinbaseResponse.Asks) == 0 || len(coinbaseResponse.Asks[0]) == 0 {
		return 0, 0, fmt.Errorf("Failed to find ask prices in coinbase response")
	}

	// Convert the prices to float64:
	// Bid represents the price someone else is willing to pay, this is our sell value
	bid, err := strconv.ParseFloat(coinbaseResponse.Bids[0][0].(string), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse bid price string from coinbase response: %w", err)
	}

	// Ask represents the price someone else is asking for, this is our buy value
	ask, err := strconv.ParseFloat(coinbaseResponse.Asks[0][0].(string), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse ask price string from coinbase response: %w", err)
	}

	return ask, bid, nil
}

// KrakenExchange implements the Exchange interface for Kraken
type KrakenExchange struct{}

// GetPrices retrieves the price for a given symbol from Kraken
func (k *KrakenExchange) GetPrices(ctx context.Context, symbol string) (float64, float64, error) {
	// Construct the Kraken API URL
	url := fmt.Sprintf("https://api.kraken.com/0/public/Depth?pair=%sUSD&count=1", symbol)

	// Create a new HTTP client with a timeout
	client := http.Client{Timeout: 10 * time.Second}

	// Send the request to the Kraken API
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get price from kraken: %w", err)
	}
	defer resp.Body.Close()

	// Define the JSON structure
	type ResultBlock struct {
		Asks [][]any `json:"asks"`
		Bids [][]any `json:"bids"`
	}

	var krakenResponse struct {
		Error []string `json:"error"`
		/*
		At first it looked like the kraken response was following
		pattern X{symbol}Z{currency}, however when calling with
		symbol BTC the response was XXBTZUSD, since we can't rely
		on knowing the key we will just use a map and grab the
		first element
		*/
		Result map[string]ResultBlock `json:"result"`
	}

	// Decode the JSON response
	if err := json.NewDecoder(resp.Body).Decode(&krakenResponse); err != nil {
		return 0, 0, fmt.Errorf("failed to decode kraken response: %w", err)
	}

	if len(krakenResponse.Error) != 0 {
		return 0, 0, fmt.Errorf("Kraken price fetch failed with errors: %s", strings.Join(krakenResponse.Error, ", "))
	}

	var ask, bid float64

	for _, aResult := range krakenResponse.Result {
		// Make sure bid and ask data are present
		if len(aResult.Bids) == 0 || len(aResult.Bids[0]) == 0 {
			return 0, 0, fmt.Errorf("Failed to find bid prices in kraken response")
		}
		if len(aResult.Asks) == 0 || len(aResult.Asks[0]) == 0 {
			return 0, 0, fmt.Errorf("Failed to find ask prices in kraken response")
		}

		bid, err = strconv.ParseFloat(aResult.Bids[0][0].(string), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse bid price string from kraken response: %w", err)
		}

		ask, err = strconv.ParseFloat(aResult.Asks[0][0].(string), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse ask price string from kraken response: %w", err)
		}
	}


	return ask, bid, nil
}

// GetName returns the name of the exchange
func (c *CoinbaseExchange) GetName() string {
	return "coinbase"
}

// GetName returns the name of the exchange
func (k *KrakenExchange) GetName() string {
	return "kraken"
}

