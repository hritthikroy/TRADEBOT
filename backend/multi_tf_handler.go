package main

import (
	"github.com/gofiber/fiber/v2"
)

// MultiTFBacktestRequest holds request parameters
type MultiTFBacktestRequest struct {
	Symbol       string  `json:"symbol"`
	Days         int     `json:"days"`
	StartBalance float64 `json:"startBalance"`
}

// HandleMultiTimeframeBacktest runs multi-timeframe backtest
func HandleMultiTimeframeBacktest(c *fiber.Ctx) error {
	var req MultiTFBacktestRequest
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
		req.Days = 7 // Multi-TF needs less days due to 1m data
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}

	// Run multi-timeframe backtest
	result, err := RunMultiTimeframeBacktest(req.Symbol, req.Days, req.StartBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"strategy":     "Multi-Timeframe Top-Down (4h→1h→15m→3m→1m)",
		"symbol":       req.Symbol,
		"days":         req.Days,
		"startBalance": req.StartBalance,
		"result":       result,
		"description":  "Uses 4h for direction, 1h for key levels, 15m for volume/delta, 3m/1m for precise entry",
	})
}
