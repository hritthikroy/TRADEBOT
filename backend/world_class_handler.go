package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"github.com/gofiber/fiber/v2"
)

// HandleWorldClassBacktest handles world-class backtest requests
func HandleWorldClassBacktest(c *fiber.Ctx) error {
	var config WorldClassBacktestConfig
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
	
	// Set advanced defaults
	if config.MonteCarloRuns == 0 {
		config.MonteCarloRuns = 1000
	}
	if config.WalkForwardPeriods == 0 {
		config.WalkForwardPeriods = 5
	}
	
	fmt.Printf("🚀 Starting world-class backtest: %s %s %dd %s\n", 
		config.Symbol, config.Interval, config.Days, config.Strategy)
	
	// Fetch candles
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
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
	
	// Run world-class backtest
	result, err := RunWorldClassBacktest(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Backtest failed: " + err.Error(),
		})
	}
	
	fmt.Printf("✅ Backtest complete: %d trades, %.2f%% WR, %.2f PF\n",
		result.TotalTrades, result.WinRate, result.ProfitFactor)
	
	return c.JSON(result)
}

// HandleWorldClassBacktestHTTP handles world-class backtest via HTTP
func HandleWorldClassBacktestHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var config WorldClassBacktestConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
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
	if config.MonteCarloRuns == 0 {
		config.MonteCarloRuns = 1000
	}
	if config.WalkForwardPeriods == 0 {
		config.WalkForwardPeriods = 5
	}
	
	// Fetch candles
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if len(candles) < 200 {
		http.Error(w, fmt.Sprintf("Insufficient data: got %d candles, need at least 200", len(candles)), http.StatusBadRequest)
		return
	}
	
	// Run world-class backtest
	result, err := RunWorldClassBacktest(config, candles)
	if err != nil {
		http.Error(w, "Backtest failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleQuickCompare compares standard vs world-class backtest
func HandleQuickCompare(c *fiber.Ctx) error {
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
	
	// Fetch candles
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}
	
	startTime := time.Now()
	
	// Run standard backtest
	standardResult, err := RunBacktest(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Standard backtest failed: " + err.Error(),
		})
	}
	standardTime := time.Since(startTime)
	
	// Run world-class backtest
	wcConfig := WorldClassBacktestConfig{
		BacktestConfig:     config,
		EnableMonteCarlo:   true,
		MonteCarloRuns:     1000,
		EnableWalkForward:  true,
		WalkForwardPeriods: 5,
		EnableStressTest:   true,
	}
	
	startTime = time.Now()
	wcResult, err := RunWorldClassBacktest(wcConfig, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "World-class backtest failed: " + err.Error(),
		})
	}
	wcTime := time.Since(startTime)
	
	return c.JSON(fiber.Map{
		"standard": fiber.Map{
			"result":   standardResult,
			"duration": standardTime.String(),
		},
		"worldClass": fiber.Map{
			"result":   wcResult,
			"duration": wcTime.String(),
		},
		"comparison": fiber.Map{
			"additionalMetrics": []string{
				"Sharpe Ratio",
				"Sortino Ratio",
				"Calmar Ratio",
				"Monte Carlo Analysis",
				"Walk Forward Analysis",
				"Stress Test Results",
				"Win/Loss Streaks",
				"Expectancy Per Trade",
				"Recovery Factor",
			},
			"recommendation": getRecommendation(wcResult),
		},
	})
}

func getRecommendation(result *WorldClassBacktestResult) string {
	if result.SharpeRatio > 2.0 && result.ProfitFactor > 2.0 && result.WinRate > 55 {
		return "⭐⭐⭐⭐⭐ Excellent - Ready for live trading"
	} else if result.SharpeRatio > 1.5 && result.ProfitFactor > 1.5 && result.WinRate > 50 {
		return "⭐⭐⭐⭐ Good - Consider paper trading first"
	} else if result.SharpeRatio > 1.0 && result.ProfitFactor > 1.2 && result.WinRate > 45 {
		return "⭐⭐⭐ Fair - Needs optimization"
	} else if result.ProfitFactor > 1.0 {
		return "⭐⭐ Poor - Significant improvements needed"
	}
	return "⭐ Very Poor - Not recommended for trading"
}
