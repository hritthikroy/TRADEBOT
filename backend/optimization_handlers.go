package main

import (
	"github.com/gofiber/fiber/v2"
)

// OptimizeParametersRequest represents optimization request
type OptimizeParametersRequest struct {
	StrategyName string  `json:"strategyName"`
	Symbol       string  `json:"symbol"`
	StartBalance float64 `json:"startBalance"`
	Days         int     `json:"days"`
}

// OptimizeAllRequest represents optimize all request
type OptimizeAllRequest struct {
	Symbol       string  `json:"symbol"`
	StartBalance float64 `json:"startBalance"`
	Days         int     `json:"days"`
}

// HandleOptimizeParameters optimizes parameters for a single strategy
func HandleOptimizeParameters(c *fiber.Ctx) error {
	var req OptimizeParametersRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}
	if req.Days == 0 {
		req.Days = 90
	}
	
	// Optimize
	results, err := OptimizeStrategyParameters(req.StrategyName, req.Symbol, req.StartBalance, req.Days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	// Return top 10 results
	topResults := results
	if len(topResults) > 10 {
		topResults = results[:10]
	}
	
	return c.JSON(fiber.Map{
		"strategyName": req.StrategyName,
		"symbol":       req.Symbol,
		"startBalance": req.StartBalance,
		"days":         req.Days,
		"totalTests":   len(results),
		"topResults":   topResults,
		"bestResult":   results[0],
	})
}

// HandleOptimizeAllStrategies optimizes all strategies
func HandleOptimizeAllStrategies(c *fiber.Ctx) error {
	var req OptimizeAllRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.StartBalance == 0 {
		req.StartBalance = 500
	}
	if req.Days == 0 {
		req.Days = 90
	}
	
	// Optimize all
	allResults, err := OptimizeAllStrategies(req.Symbol, req.StartBalance, req.Days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	// Collect best result for each strategy
	bestResults := make(map[string]OptimizationResult)
	for strategyName, results := range allResults {
		if len(results) > 0 {
			bestResults[strategyName] = results[0]
		}
	}
	
	// Find overall best
	var overallBest OptimizationResult
	bestScore := 0.0
	for _, result := range bestResults {
		if result.Score > bestScore {
			bestScore = result.Score
			overallBest = result
		}
	}
	
	return c.JSON(fiber.Map{
		"symbol":       req.Symbol,
		"startBalance": req.StartBalance,
		"days":         req.Days,
		"totalStrategies": len(bestResults),
		"bestResults": bestResults,
		"overallBest": overallBest,
	})
}
