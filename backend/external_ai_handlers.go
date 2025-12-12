package main

import (
	"fmt"
	
	"github.com/gofiber/fiber/v2"
)

// HandleEnhancedAIAnalysis uses external AI APIs for enhanced analysis
func HandleEnhancedAIAnalysis(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	interval := c.Query("interval", "15m")
	days := c.QueryInt("days", 30)
	
	fmt.Printf("🤖 Enhanced AI analysis with external APIs: %s %s %dd\n", symbol, interval, days)
	
	// Fetch candles
	candles, err := FetchHistoricalData(symbol, interval, days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}
	
	if len(candles) < 200 {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Insufficient data: got %d candles, need at least 200", len(candles)),
		})
	}
	
	// Run local AI analysis first
	localAnalysis, err := AnalyzeMarketWithAI(candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Local AI analysis failed: " + err.Error(),
		})
	}
	
	// Enhance with external AI
	enhanced, err := EnhanceWithExternalAI(localAnalysis, candles, symbol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "External AI enhancement failed: " + err.Error(),
		})
	}
	
	fmt.Printf("✅ Enhanced AI analysis complete\n")
	fmt.Printf("   Combined Recommendation: %s\n", enhanced.CombinedRecommendation)
	fmt.Printf("   Confidence: %.1f%%\n", enhanced.ConfidenceScore)
	
	return c.JSON(enhanced)
}

// HandleAIConfig returns current AI configuration status
func HandleAIConfig(c *fiber.Ctx) error {
	config := LoadAIConfig()
	
	return c.JSON(fiber.Map{
		"openai": fiber.Map{
			"enabled": config.EnableOpenAI,
			"configured": config.OpenAIKey != "",
		},
		"anthropic": fiber.Map{
			"enabled": config.EnableAnthropic,
			"configured": config.AnthropicKey != "",
		},
		"gemini": fiber.Map{
			"enabled": config.EnableGemini,
			"configured": config.GeminiKey != "",
		},
		"perplexity": fiber.Map{
			"enabled": config.EnablePerplexity,
			"configured": config.PerplexityKey != "",
		},
		"message": "Set API keys in environment variables to enable external AI",
	})
}
