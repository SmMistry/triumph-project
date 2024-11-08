package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/SmMistry/triumph-project/services/order"
	"github.com/SmMistry/triumph-project/controllers/orders"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// MockExchange is a mock implementation of the Exchange interface for testing.
type MockExchange struct {
	Name  string
	BuyPrice float64
	SellPrice float64
	Err   error
}

func (m *MockExchange) GetPrices(ctx context.Context, symbol string) (float64, float64, error) {
	return m.BuyPrice, m.SellPrice, m.Err
}

func (m *MockExchange) GetName() string {
	return m.Name
}

func TestBuyHandler(t *testing.T) {
	tests := []struct {
		name            string
		amount          string
		symbol          string
		mockExchanges   []*MockExchange
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name: "Valid request with best price on Coinbase",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 9900, SellPrice: 9900, Err: nil},
				{Name: "kraken", BuyPrice: 10000, SellPrice: 10000, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["coinbase"],"usdAmount":9900}`,
		},
		{
			name: "Valid request for ETH with best price on Coinbase",
			amount: "1",
			symbol: "ETH",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 9900, SellPrice: 9900, Err: nil},
				{Name: "kraken", BuyPrice: 10000, SellPrice: 10000, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"ETH","exchange":["coinbase"],"usdAmount":9900}`,
		},
		{
			name: "Valid request with best price on Kraken",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["kraken"],"usdAmount":9900}`,
		},
		{
			name: "Valid request with same price on both exchanges",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 10000, SellPrice: 10000, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["coinbase","kraken"],"usdAmount":10000}`,
		},
		{
			name: "Valid request with fractional amount best price on Kraken",
			amount: "0.5",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":0.5,"coin":"BTC","exchange":["kraken"],"usdAmount":4950}`,
		},
		{
			name: "Invalid amount parameter",
			amount: "junk",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid amount"}`,
		},
		{
			name: "Missing amount parameter",
			amount: "",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid amount"}`,
		},
		{
			name: "Error fetching price from Coinbase",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("coinbase error")},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["kraken"],"usdAmount":9900}`,
		},
		{
			name: "Error fetching price from both exchanges",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("coinbase error")},
				{Name: "kraken", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("kraken error")},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"failed to find best price for BTC"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Fiber app
			app := fiber.New()

			// Create a new OrderService with mock exchanges
			coinbase := tt.mockExchanges[0]
			kraken := tt.mockExchanges[1]
			orderService := order.NewOrderService(coinbase, kraken)

			// Create a new OrderController
			orderController := orders.NewOrderController(orderService)

			// Define the API route
			app.Get("/buy", orderController.BuyHandler)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/buy?amount=%s&symbol=%s", tt.amount, tt.symbol), nil)

			// Perform the request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Assert the response body
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, string(body))

			// Assert the response headers (if any)
			if tt.expectedHeaders != nil {
				for key, value := range tt.expectedHeaders {
					assert.Equal(t, value, resp.Header.Get(key))
				}
			}
		})
	}
}

func TestSellHandler(t *testing.T) {
	tests := []struct {
		name            string
		amount          string
		symbol          string
		mockExchanges   []*MockExchange
		expectedStatus  int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name: "Valid request with best price on Coinbase",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["coinbase"],"usdAmount":10000}`,
		},
		{
			name: "Valid request for ETH with best price on Coinbase",
			amount: "1",
			symbol: "ETH",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"ETH","exchange":["coinbase"],"usdAmount":10000}`,
		},
		{
			name: "Valid request with best price on Kraken",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 9900, SellPrice: 9900, Err: nil},
				{Name: "kraken", BuyPrice: 10000, SellPrice: 10000, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["kraken"],"usdAmount":10000}`,
		},
		{
			name: "Valid request with same price on both exchanges",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 9900, SellPrice: 9900, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["coinbase", "kraken"],"usdAmount":9900}`,
		},
		{
			name: "Valid request with fractional amount and best price on Kraken",
			amount: "0.5",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 9900, SellPrice: 9900, Err: nil},
				{Name: "kraken", BuyPrice: 10000, SellPrice: 10000, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":0.5,"coin":"BTC","exchange":["kraken"],"usdAmount":5000}`,
		},
		{
			name: "Invalid amount parameter",
			amount: "invalid",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid amount"}`,
		},
		{
			name: "Missing amount parameter",
			amount: "",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 10000, SellPrice: 10000, Err: nil},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"error":"invalid amount"}`,
		},
		{
			name: "Error fetching price from Coinbase",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("coinbase error")},
				{Name: "kraken", BuyPrice: 9900, SellPrice: 9900, Err: nil},
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"amount":1,"coin":"BTC","exchange":["kraken"],"usdAmount":9900}`,
		},
		{
			name: "Error fetching price from both exchanges",
			amount: "1",
			symbol: "BTC",
			mockExchanges: []*MockExchange{
				{Name: "coinbase", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("coinbase error")},
				{Name: "kraken", BuyPrice: 0, SellPrice: 0, Err: fmt.Errorf("kraken error")},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{"error":"failed to find best price for BTC"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Fiber app
			app := fiber.New()

			// Create a new OrderService with mock exchanges
			coinbase := tt.mockExchanges[0]
			kraken := tt.mockExchanges[1]
			orderService := order.NewOrderService(coinbase, kraken)

			// Create a new OrderController
			orderController := orders.NewOrderController(orderService)

			// Define the API route
			app.Get("/sell", orderController.SellHandler)

			// Create a test request
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sell?amount=%s&symbol=%s", tt.amount, tt.symbol), nil)

			// Perform the request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Assert the response status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Assert the response body
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedBody, string(body))

			// Assert the response headers (if any)
			if tt.expectedHeaders != nil {
				for key, value := range tt.expectedHeaders {
					assert.Equal(t, value, resp.Header.Get(key))
				}
			}
		})
	}
}