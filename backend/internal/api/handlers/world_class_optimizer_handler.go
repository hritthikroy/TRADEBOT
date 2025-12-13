package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// HandleWorldClassOptimization runs world-class parameter optimization
func HandleWorldClassOptimization(c *fiber.Ctx) error {
	log.Println("🌍 Starting World-Class Optimization...")
	
	// Create optimizer
	optimizer := NewWorldClassOptimizer()
	
	// Run optimization (this will take a while!)
	results := optimizer.OptimizeAll()
	
	// Save results
	filename := "WORLD_CLASS_OPTIMIZATION_RESULTS.json"
	err := optimizer.SaveResults(results, filename)
	if err != nil {
		log.Printf("❌ Failed to save results: %v", err)
	} else {
		log.Printf("✅ Results saved to: %s", filename)
	}
	
	// Print summary
	results.PrintSummary()
	
	return c.JSON(results)
}

// HandleQuickOptimization runs a faster optimization with fewer parameters
func HandleQuickOptimization(c *fiber.Ctx) error {
	var req struct {
		Strategy string `json:"strategy"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	
	log.Printf("⚡ Quick optimization for: %s", req.Strategy)
	
	// Fetch candles
	candles, err := fetchBinanceData("BTCUSDT", "15m", 180)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data",
		})
	}
	
	// Create optimizer
	optimizer := NewWorldClassOptimizer()
	
	// Optimize single strategy
	result := optimizer.OptimizeStrategy(req.Strategy, candles)
	
	return c.JSON(result)
}
