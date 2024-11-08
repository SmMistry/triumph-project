package orders

import (
	"net/http"

	"github.com/SmMistry/triumph-project/services/order"
	"github.com/gofiber/fiber/v2"
)

// OrderController handles HTTP requests for orders
type OrderController struct {
	orderService *order.OrderService
}

// NewOrderController creates a new OrderController with the given OrderService
func NewOrderController(orderService *order.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

// BuyHandler handles the /buy endpoint
func (oc *OrderController) BuyHandler(c *fiber.Ctx) error {
	// Parse the request parameters
	amount := c.QueryFloat("amount", 0)
	if amount == 0{
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid amount"})
	}
	symbol := c.Query("symbol")

	// Execute the buy order
	usdAmount, exchanges, err := oc.orderService.Buy(c.Context(), amount, symbol)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the response
	return c.JSON(fiber.Map{
		"coin": symbol,
		"amount": amount,
		"usdAmount": usdAmount,
		"exchange":  exchanges,
	})
}

// SellHandler
func (oc *OrderController) SellHandler(c *fiber.Ctx) error {
	// Parse the request parameters
	amount := c.QueryFloat("amount", 0)
	if amount == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid amount"})
	}
	symbol := c.Query("symbol")

	// Execute the sell order
	usdAmount, exchanges, err := oc.orderService.Sell(c.Context(), amount, symbol)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the response
	return c.JSON(fiber.Map{
		"coin": symbol,
		"amount": amount,
		"usdAmount": usdAmount,
		"exchange":  exchanges,
	})
}