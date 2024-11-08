package main

import (
	"log"
	
	"github.com/SmMistry/triumph-project/services/exchange"
	"github.com/SmMistry/triumph-project/services/order"
	"github.com/SmMistry/triumph-project/controllers/orders"

	"github.com/gofiber/fiber/v2"
)

func initializeService() *order.OrderService {
	// Initialize the exchanges
	coinbase := &exchange.CoinbaseExchange{}
	kraken := &exchange.KrakenExchange{}

	return order.NewOrderService(coinbase, kraken)
}

func initializeOrderController(orderService *order.OrderService) *orders.OrderController{
	return orders.NewOrderController(orderService)
}

func main() {
	// Create the order service
	orderService := initializeService()

	// Create the order controller
	orderController := initializeOrderController(orderService)

	// Initialize the Fiber app
	app := fiber.New()

	// Define the API routes
	app.Get("/buy", orderController.BuyHandler)
	app.Get("/sell", orderController.SellHandler)

	// Start the server
	log.Fatal(app.Listen(":4000"))
}