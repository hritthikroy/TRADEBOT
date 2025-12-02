package main

import (
	"github.com/gofiber/fiber/v2"
)

// GetAIStats returns AI enhancement statistics
func GetAIStats(c *fiber.Ctx) error {
	aiGen := c.Locals("aiGenerator").(*AIEnhancedSignalGenerator)
	
	stats := aiGen.GetAIStats()
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}

// ToggleAIFilter enables or disables AI filtering
func ToggleAIFilter(c *fiber.Ctx) error {
	aiGen := c.Locals("aiGenerator").(*AIEnhancedSignalGenerator)
	
	type ToggleRequest struct {
		Enabled bool `json:"enabled"`
	}
	
	var req ToggleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	
	aiGen.ToggleAI(req.Enabled)
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "AI filter toggled",
		"enabled": req.Enabled,
	})
}

// TestAIConnection tests the Grok AI connection
func TestAIConnection(c *fiber.Ctx) error {
	aiGen := c.Locals("aiGenerator").(*AIEnhancedSignalGenerator)
	
	// Test with a simple query
	response, err := aiGen.GrokService.CallGrokAPI("Say 'AI connection successful' if you can read this.")
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Grok AI connection successful",
		"response": response,
	})
}

// AnalyzeSymbolSentiment analyzes sentiment for a specific symbol
func AnalyzeSymbolSentiment(c *fiber.Ctx) error {
	aiGen := c.Locals("aiGenerator").(*AIEnhancedSignalGenerator)
	
	type SentimentRequest struct {
		Symbol       string  `json:"symbol"`
		SignalType   string  `json:"signal_type"`
		CurrentPrice float64 `json:"current_price"`
		Strength     int     `json:"strength"`
	}
	
	var req SentimentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	
	// Default values
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.SignalType == "" {
		req.SignalType = "BUY"
	}
	if req.CurrentPrice == 0 {
		req.CurrentPrice = 50000.0
	}
	if req.Strength == 0 {
		req.Strength = 70
	}
	
	sentiment, err := aiGen.GrokService.AnalyzeMarketSentiment(
		req.Symbol,
		req.SignalType,
		req.CurrentPrice,
		req.Strength,
	)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"data":    sentiment,
	})
}
