package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// AllTimeframesRequest holds request parameters
type AllTimeframesRequest struct {
	Symbol       string  `json:"symbol"`
	StartBalance float64 `json:"startBalance"`
}

// HandleAllTimeframesBacktest tests all timeframes with optimized strategies
func HandleAllTimeframesBacktest(c *fiber.Ctx) error {
	var req AllTimeframesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}

	// Run all timeframes backtest
	results, err := RunAllTimeframesBacktest(req.Symbol, req.StartBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"symbol":       req.Symbol,
		"startBalance": req.StartBalance,
		"results":      results,
	})
}
