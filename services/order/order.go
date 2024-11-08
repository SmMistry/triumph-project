package order

import (
	"fmt"
	"log"
	"context"
	"triumph/services/exchange"
)

// OrderService handles order execution logic
type OrderService struct {
	exchanges []exchange.Exchange
}

// NewOrderService creates a new OrderService with the given exchanges
func NewOrderService(exchanges ...exchange.Exchange) *OrderService {
	return &OrderService{exchanges: exchanges}
}

// Buy executes a buy order for the given amount and symbol
func (o *OrderService) Buy(ctx context.Context, amount float64, symbol string) (float64, []string, error) {
	bestPrice := 1e18 // Initialize with a very high value
	bestExchanges := []string{}

	// Iterate over the exchanges to find the best price
	for _, exchange := range o.exchanges {
		price, _, err := exchange.GetPrices(ctx, symbol)
		if err != nil {
			log.Printf("failed to get price from exchange: %v", err)
			continue
		}
		// log.Printf("Found price: %v for exchange %s", price, exchange.GetName())

		if price < bestPrice {
			bestPrice = price
			bestExchanges = []string{exchange.GetName()}
		} else if price == bestPrice {
			bestExchanges = append(bestExchanges, exchange.GetName())
		}
	}

	// If no best price was found, return an error
	if bestPrice == 1e18 {
		return 0, nil, fmt.Errorf("failed to find best price for %s", symbol)
	}

	// Calculate the USD amount
	usdAmount := amount * bestPrice

	return usdAmount, bestExchanges, nil
}

// Sell executes a sell order for the given amount and symbol
func (o *OrderService) Sell(ctx context.Context, amount float64, symbol string) (float64, []string, error) {
	bestPrice := 0.0 
	bestExchanges := []string{}

	// Iterate over the exchanges to find the best price
	for _, exchange := range o.exchanges {
		_, price, err := exchange.GetPrices(ctx, symbol)
		if err != nil {
			log.Printf("failed to get price from exchange: %v", err)
			continue
		}

		if price > bestPrice {
			bestPrice = price
			bestExchanges = []string{exchange.GetName()}
		} else if price == bestPrice {
			bestExchanges = append(bestExchanges, exchange.GetName())
		}
	}

	// If no best price was found, return an error
	if bestPrice == 0.0 {
		return 0, nil, fmt.Errorf("failed to find best price for %s", symbol)
	}

	// Calculate the USD amount
	usdAmount := amount * bestPrice

	return usdAmount, bestExchanges, nil
}
