package main

import (
	"github.com/gofiber/fiber/v2"
)

// HandleTestAllStrategies tests all advanced strategies
func HandleTestAllStrategies(c *fiber.Ctx) error {
	var req struct {
		Symbol       string  `json:"symbol"`
		StartBalance float64 `json:"startBalance"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}
	
	// Test all strategies
	results, err := TestAllStrategies(req.Symbol, req.StartBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"symbol":         req.Symbol,
		"startBalance":   req.StartBalance,
		"totalStrategies": len(results),
		"results":        results,
		"bestStrategy":   results[0],
	})
}
