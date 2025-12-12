package main

import (
	"fmt"
	"time"
	
	"github.com/gofiber/fiber/v2"
)

// HandleAIOptimization runs AI-powered strategy optimization
func HandleAIOptimization(c *fiber.Ctx) error {
	var config BacktestConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request: " + err.Error(),
		})
	}
	
	// Set defaults
	if config.Symbol == "" {
		config.Symbol = "BTCUSDT"
	}
	if config.Interval == "" {
		config.Interval = "15m"
	}
	if config.Days == 0 {
		config.Days = 30
	}
	if config.StartBalance == 0 {
		config.StartBalance = 1000
	}
	if config.Strategy == "" {
		config.Strategy = "session_trader"
	}
	
	fmt.Printf("🤖 Starting AI optimization: %s %s %dd %s\n", 
		config.Symbol, config.Interval, config.Days, config.Strategy)
	
	// Fetch candles
	candles, err := FetchHistoricalData(config.Symbol, config.Interval, config.Days)
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
	
	startTime := time.Now()
	
	// Run AI optimization
	result, err := RunAIOptimization(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "AI optimization failed: " + err.Error(),
		})
	}
	
	result.Duration = time.Since(startTime).String()
	
	fmt.Printf("✅ AI optimization complete in %v\n", time.Since(startTime))
	fmt.Printf("   Best fitness: %.2f, Improvement: %.2f%%\n", 
		result.BestStrategy.Fitness, result.ImprovementPct)
	
	return c.JSON(result)
}

// HandleAIMarketAnalysis provides AI-powered market insights
func HandleAIMarketAnalysis(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	interval := c.Query("interval", "15m")
	days := c.QueryInt("days", 30)
	
	fmt.Printf("🤖 AI analyzing market: %s %s %dd\n", symbol, interval, days)
	
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
	
	// Run AI analysis
	analysis, err := AnalyzeMarketWithAI(candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "AI analysis failed: " + err.Error(),
		})
	}
	
	fmt.Printf("✅ AI analysis complete\n")
	fmt.Printf("   Market: %s, Trend: %s, Volatility: %s\n", 
		analysis.MarketRegime, analysis.PredictedMove, analysis.VolatilityLevel)
	
	return c.JSON(analysis)
}

// HandleAIStrategyRecommendation recommends best strategy for current market
func HandleAIStrategyRecommendation(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	interval := c.Query("interval", "15m")
	days := c.QueryInt("days", 30)
	
	// Fetch candles
	candles, err := FetchHistoricalData(symbol, interval, days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}
	
	// Analyze market
	analysis, err := AnalyzeMarketWithAI(candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "AI analysis failed: " + err.Error(),
		})
	}
	
	// Test recommended strategy
	config := BacktestConfig{
		Symbol:       symbol,
		Interval:     interval,
		Days:         days,
		StartBalance: 1000,
		Strategy:     analysis.BestStrategy,
	}
	
	result, err := RunBacktest(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Backtest failed: " + err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"marketAnalysis": analysis,
		"recommendedStrategy": analysis.BestStrategy,
		"backtestResult": fiber.Map{
			"totalTrades":  result.TotalTrades,
			"winRate":      result.WinRate,
			"profitFactor": result.ProfitFactor,
			"maxDrawdown":  result.MaxDrawdown,
			"returnPercent": result.ReturnPercent,
		},
		"confidence": analysis.Confidence,
		"riskLevel": analysis.RiskLevel,
	})
}

// HandleAICompareStrategies uses AI to compare all strategies for current market
func HandleAICompareStrategies(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	interval := c.Query("interval", "15m")
	days := c.QueryInt("days", 30)
	
	fmt.Printf("🤖 AI comparing all strategies: %s %s %dd\n", symbol, interval, days)
	
	// Fetch candles
	candles, err := FetchHistoricalData(symbol, interval, days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}
	
	// Analyze market first
	analysis, err := AnalyzeMarketWithAI(candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "AI analysis failed: " + err.Error(),
		})
	}
	
	strategies := []string{
		"session_trader",
		"liquidity_hunter",
		"breakout_master",
		"trend_rider",
		"range_master",
		"smart_money_tracker",
		"institutional_follower",
		"reversal_sniper",
		"momentum_beast",
		"scalper_pro",
	}
	
	results := make([]fiber.Map, 0)
	
	for _, strategy := range strategies {
		config := BacktestConfig{
			Symbol:       symbol,
			Interval:     interval,
			Days:         days,
			StartBalance: 1000,
			Strategy:     strategy,
		}
		
		result, err := RunBacktest(config, candles)
		if err != nil {
			continue
		}
		
		// Calculate AI score
		aiScore := calculateAIScore(result, analysis)
		
		results = append(results, fiber.Map{
			"strategy":     strategy,
			"totalTrades":  result.TotalTrades,
			"winRate":      result.WinRate,
			"profitFactor": result.ProfitFactor,
			"maxDrawdown":  result.MaxDrawdown,
			"returnPercent": result.ReturnPercent,
			"aiScore":      aiScore,
			"recommended":  strategy == analysis.BestStrategy,
		})
	}
	
	return c.JSON(fiber.Map{
		"marketAnalysis": analysis,
		"strategies": results,
		"recommendation": fmt.Sprintf("Based on AI analysis, %s is best for current %s market with %s volatility",
			analysis.BestStrategy, analysis.MarketRegime, analysis.VolatilityLevel),
	})
}

// calculateAIScore calculates AI-based strategy score
func calculateAIScore(result *BacktestResult, analysis *AIMarketAnalysis) float64 {
	if result.TotalTrades == 0 {
		return 0
	}
	
	// Base score from performance
	baseScore := (result.WinRate * 0.3) + (result.ProfitFactor * 20 * 0.4) + ((100 - result.MaxDrawdown) * 0.3)
	
	// Adjust based on market conditions
	if analysis.VolatilityLevel == "high" && result.MaxDrawdown < 15 {
		baseScore += 10 // Bonus for low drawdown in high volatility
	}
	if analysis.MarketRegime == "trending" && result.ProfitFactor > 2.0 {
		baseScore += 10 // Bonus for good performance in trending market
	}
	if analysis.MarketRegime == "ranging" && result.WinRate > 60 {
		baseScore += 10 // Bonus for high win rate in ranging market
	}
	
	return math.Min(baseScore, 100)
}
