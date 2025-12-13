package main

import (
	"github.com/gofiber/fiber/v2"
)

// HandleTestAllStrategies tests all advanced strategies
func HandleTestAllStrategies(c *fiber.Ctx) error {
	var req struct {
		Symbol       string  `json:"symbol"`
		Days         int     `json:"days"`         // FIX: Added days parameter
		StartBalance float64 `json:"startBalance"`
		FilterBuy    *bool   `json:"filterBuy"`    // nil = both, true = buy only, false = exclude buy
		FilterSell   *bool   `json:"filterSell"`   // nil = both, true = sell only, false = exclude sell
		StartTime    *int64  `json:"startTime"`    // Optional: Unix timestamp in milliseconds
		EndTime      *int64  `json:"endTime"`      // Optional: Unix timestamp in milliseconds
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
	
	// Default to both if not specified
	filterBuy := true
	filterSell := true
	if req.FilterBuy != nil {
		filterBuy = *req.FilterBuy
	}
	if req.FilterSell != nil {
		filterSell = *req.FilterSell
	}
	
	// Default to 30 days if not specified
	days := req.Days
	if days == 0 {
		days = 30
	}
	
	// Test all strategies with filters and optional date range
	results, err := TestAllStrategiesWithFilterAndRange(req.Symbol, days, req.StartBalance, filterBuy, filterSell, req.StartTime, req.EndTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"symbol":         req.Symbol,
		"startBalance":   req.StartBalance,
		"filterBuy":      filterBuy,
		"filterSell":     filterSell,
		"totalStrategies": len(results),
		"results":        results,
		"bestStrategy":   results[0],
	})
}
