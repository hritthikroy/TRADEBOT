package main

import (
	"github.com/gofiber/fiber/v2"
)

// ComprehensiveBacktestRequest holds request parameters
type ComprehensiveBacktestRequest struct {
	Symbol       string  `json:"symbol"`
	Days         int     `json:"days"`
	StartBalance float64 `json:"startBalance"`
}

// HandleComprehensiveBacktest runs backtest across all timeframes
func HandleComprehensiveBacktest(c *fiber.Ctx) error {
	var req ComprehensiveBacktestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}

	// Run comprehensive backtest
	result, err := RunComprehensiveBacktest(req.Symbol, req.Days, req.StartBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
