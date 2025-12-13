package handlers

import (
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HandleOptimizedBacktest runs backtest with optimized daily strategies
func HandleOptimizedBacktest(c *fiber.Ctx) error {
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

	fmt.Printf("🚀 Running OPTIMIZED backtest: %s %s %dd %s\n",
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

	startTime := time.Now()

	// Run optimized backtest
	result, err := RunOptimizedBacktest(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Backtest failed: " + err.Error(),
		})
	}

	result.Duration = time.Since(startTime).String()

	fmt.Printf("✅ Optimized backtest complete in %v\n", time.Since(startTime))
	fmt.Printf("   Trades: %d, WinRate: %.1f%%, PF: %.2f, Return: %.2f%%\n",
		result.TotalTrades, result.WinRate, result.ProfitFactor, result.ReturnPercent)

	return c.JSON(result)
}

// HandleOptimizeAllDailyStrategies tests all 10 optimized strategies
func HandleOptimizeAllDailyStrategies(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "BTCUSDT")
	interval := c.Query("interval", "15m")
	days := c.QueryInt("days", 30)
	startBalance := c.QueryFloat("startBalance", 1000)

	fmt.Printf("🚀 Testing ALL 10 optimized daily strategies: %s %s %dd\n", symbol, interval, days)

	// Fetch candles
	candles, err := fetchBinanceData(symbol, interval, days)
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

	strategies := GetOptimizedDailyStrategies()
	results := make([]fiber.Map, 0)

	for strategyName, strategy := range strategies {
		config := BacktestConfig{
			Symbol:       symbol,
			Interval:     interval,
			Days:         days,
			StartBalance: startBalance,
			Strategy:     strategyName,
		}

		result, err := RunOptimizedBacktest(config, candles)
		if err != nil {
			continue
		}

		results = append(results, fiber.Map{
			"strategy":        strategyName,
			"name":            strategy.Name,
			"description":     strategy.Description,
			"timeframe":       strategy.Timeframe,
			"totalTrades":     result.TotalTrades,
			"winningTrades":   result.WinningTrades,
			"losingTrades":    result.LosingTrades,
			"winRate":         result.WinRate,
			"profitFactor":    result.ProfitFactor,
			"maxDrawdown":     result.MaxDrawdown,
			"returnPercent":   result.ReturnPercent,
			"finalBalance":    result.FinalBalance,
			"targetWinRate":   strategy.TargetWinRate,
			"targetPF":        strategy.TargetProfitFactor,
			"optimizedFor":    strategy.OptimizedFor,
			"maxDailyTrades":  strategy.MaxDailyTrades,
			"riskReward":      strategy.RiskRewardRatio,
		})

		fmt.Printf("✅ %s: %d trades, %.1f%% WR, %.2f PF, %.2f%% return\n",
			strategyName, result.TotalTrades, result.WinRate, result.ProfitFactor, result.ReturnPercent)
	}

	return c.JSON(fiber.Map{
		"symbol":     symbol,
		"interval":   interval,
		"days":       days,
		"strategies": results,
		"summary": fiber.Map{
			"totalStrategies": len(results),
			"message":         "All 10 strategies optimized for daily trading profitability",
		},
	})
}

// RunOptimizedBacktest runs backtest with optimized strategy parameters
func RunOptimizedBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		PeakBalance:  config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}

	// Get optimized strategy
	optimized := OptimizeDailyStrategyParameters(candles, config.Strategy)
	if optimized == nil {
		return nil, fmt.Errorf("strategy not found: %s", config.Strategy)
	}

	// Set defaults
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.01 // 1% risk per trade
	}

	skipAhead := 10
	tradesThisDay := 0
	currentDay := ""

	// Simulate trading
	for i := 100; i < len(candles)-10; i++ {
		// Reset daily trade counter
		day := time.Unix(candles[i].Timestamp/1000, 0).Format("2006-01-02")
		if day != currentDay {
			currentDay = day
			tradesThisDay = 0
		}

		// Check max daily trades
		if tradesThisDay >= optimized.MaxDailyTrades {
			continue
		}

		// Check trading hours
		if !isInTradingHours(candles[i].Timestamp, optimized.TradingHours) {
			continue
		}

		// Generate signal using optimized strategy
		signal := GenerateOptimizedSignal(candles[:i+1], config.Strategy)
		if signal == nil {
			continue
		}

		// Simulate trade
		futureData := candles[i : minInt(i+50, len(candles))]
		trade := simulateTrade(signal, futureData, result.FinalBalance, config)

		if trade != nil {
			trade.EntryIndex = i
			trade.BalanceAfter = result.FinalBalance + trade.Profit

			result.Trades = append(result.Trades, *trade)
			result.TotalTrades++
			tradesThisDay++

			if trade.Profit > 0 {
				result.WinningTrades++
				result.TotalProfit += trade.Profit
			} else {
				result.LosingTrades++
				result.TotalLoss += math.Abs(trade.Profit)
			}

			// Update balance
			result.FinalBalance += trade.Profit

			// Track peak and drawdown
			if result.FinalBalance > result.PeakBalance {
				result.PeakBalance = result.FinalBalance
			}

			drawdown := (result.PeakBalance - result.FinalBalance) / result.PeakBalance
			if drawdown > result.MaxDrawdown {
				result.MaxDrawdown = drawdown
			}

			// Track exit reasons
			result.ExitReasons[trade.ExitReason]++

			// Skip ahead
			i += skipAhead
		}
	}

	// Calculate statistics
	calculateStats(result)

	return result, nil
}
