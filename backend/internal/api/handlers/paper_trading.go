package handlers

import (
	"time"
	
	"github.com/gofiber/fiber/v2"
)

// SetupPaperTradingRoutes sets up paper trading API routes
func SetupPaperTradingRoutes(app *fiber.App) {
	api := app.Group("/api/v1/paper-trading")
	
	// Get paper trading statistics
	api.Get("/stats", getPaperTradingStats)
	
	// Get all paper trades
	api.Get("/trades", getAllPaperTrades)
	
	// Add a new paper trade manually
	api.Post("/trade", addPaperTrade)
	
	// Update open trades with current price
	api.Post("/update", updatePaperTrades)
	
	// Reset paper trading
	api.Post("/reset", resetPaperTrading)
	
	// Start auto paper trading
	api.Post("/start-auto", startAutoPaperTrading)
	
	// Stop auto paper trading
	api.Post("/stop-auto", stopAutoPaperTrading)
	
	// Start/stop individual strategies
	api.Post("/start", startStrategyPaperTrading)
	api.Post("/stop", stopStrategyPaperTrading)
}

// getPaperTradingStats returns paper trading statistics
func getPaperTradingStats(c *fiber.Ctx) error {
	stats := paperTradingManager.GetStats()
	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// getAllPaperTrades returns all paper trades
func getAllPaperTrades(c *fiber.Ctx) error {
	trades := paperTradingManager.GetAllTrades()
	return c.JSON(fiber.Map{
		"success": true,
		"trades":  trades,
		"total":   len(trades),
	})
}

// addPaperTrade adds a new paper trade
func addPaperTrade(c *fiber.Ctx) error {
	var req struct {
		Symbol string `json:"symbol"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}
	
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	
	// Get current signal (need ~300 candles for indicators)
	candles, err := fetchBinanceData(req.Symbol, "15m", 4) // 4 days = ~384 candles
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch candles",
		})
	}
	
	generator := &UnifiedSignalGenerator{}
	signal := generator.GenerateSignal(candles, "session_trader")
	
	if signal == nil || signal.Type == "NONE" {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "No signal available",
		})
	}
	
	currentPrice := candles[len(candles)-1].Close
	trade := paperTradingManager.AddTrade(signal, currentPrice)
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Paper trade added",
		"trade":   trade,
	})
}

// updatePaperTrades updates all open trades
func updatePaperTrades(c *fiber.Ctx) error {
	var req struct {
		Symbol string `json:"symbol"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}
	
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	
	// Get current price
	candles, err := fetchBinanceData(req.Symbol, "15m", 1) // Just need latest candle
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch price",
		})
	}
	
	currentPrice := candles[0].Close
	closedTrades := paperTradingManager.UpdateOpenTrades(currentPrice)
	
	return c.JSON(fiber.Map{
		"success":      true,
		"currentPrice": currentPrice,
		"closedTrades": closedTrades,
		"message":      "Trades updated",
	})
}

// resetPaperTrading resets all paper trading data
func resetPaperTrading(c *fiber.Ctx) error {
	paperTradingManager.ResetPaperTrading()
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Paper trading reset successfully",
	})
}

// Auto paper trading variables
var (
	autoPaperTradingRunning = false
	autoPaperTradingStop    chan bool
)

// startAutoPaperTrading starts automatic paper trading
func startAutoPaperTrading(c *fiber.Ctx) error {
	if autoPaperTradingRunning {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Auto paper trading already running",
		})
	}
	
	autoPaperTradingRunning = true
	autoPaperTradingStop = make(chan bool)
	
	go runAutoPaperTrading()
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto paper trading started",
	})
}

// stopAutoPaperTrading stops automatic paper trading
func stopAutoPaperTrading(c *fiber.Ctx) error {
	if !autoPaperTradingRunning {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Auto paper trading not running",
		})
	}
	
	autoPaperTradingStop <- true
	autoPaperTradingRunning = false
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Auto paper trading stopped",
	})
}

// runAutoPaperTrading runs automatic paper trading
func runAutoPaperTrading() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-autoPaperTradingStop:
			return
		case <-ticker.C:
			// Check for new signals
			candles, err := fetchBinanceData("BTCUSDT", "15m", 4) // 4 days for indicators
			if err != nil {
				continue
			}
			
			generator := &UnifiedSignalGenerator{}
			signal := generator.GenerateSignal(candles, "session_trader")
			
			if signal != nil && signal.Type != "NONE" {
				currentPrice := candles[len(candles)-1].Close
				paperTradingManager.AddTrade(signal, currentPrice)
			}
			
			// Update open trades
			currentPrice := candles[len(candles)-1].Close
			paperTradingManager.UpdateOpenTrades(currentPrice)
		}
	}
}

// addTestTrade adds a test trade for demonstration
func addTestTrade(c *fiber.Ctx) error {
	var req struct {
		Type string `json:"type"` // "winning" or "losing"
	}
	
	if err := c.BodyParser(&req); err != nil {
		req.Type = "winning" // Default to winning trade
	}
	
	// Get current price
	candles, err := fetchBinanceData("BTCUSDT", "15m", 1)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch price",
		})
	}
	
	currentPrice := candles[0].Close
	
	// Create test signal
	var signal *AdvancedSignal
	
	if req.Type == "losing" {
		// Create a losing trade (will hit stop loss)
		signal = &AdvancedSignal{
			Type:     "BUY",
			Entry:    currentPrice,
			StopLoss: currentPrice * 1.01, // SL above entry (will hit immediately)
			TP1:      currentPrice * 0.97,
			TP2:      currentPrice * 0.95,
			TP3:      currentPrice * 0.93,
		}
	} else {
		// Create a winning trade (will hit TP)
		signal = &AdvancedSignal{
			Type:     "BUY",
			Entry:    currentPrice,
			StopLoss: currentPrice * 0.997,
			TP1:      currentPrice * 1.003,
			TP2:      currentPrice * 1.006,
			TP3:      currentPrice * 1.01,
		}
	}
	
	// Add trade
	trade := paperTradingManager.AddTrade(signal, currentPrice)
	
	// Immediately update to close the trade
	if req.Type == "losing" {
		// Simulate hitting stop loss
		paperTradingManager.UpdateOpenTrades(currentPrice * 1.011)
	} else {
		// Simulate hitting TP2
		paperTradingManager.UpdateOpenTrades(currentPrice * 1.006)
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Test trade added and closed",
		"trade":   trade,
		"type":    req.Type,
	})
}

// Strategy-specific paper trading
var (
	runningStrategies = make(map[string]bool)
	strategyStopChans = make(map[string]chan bool)
)

// startStrategyPaperTrading starts paper trading for a specific strategy
func startStrategyPaperTrading(c *fiber.Ctx) error {
	var req struct {
		Strategy string `json:"strategy"`
		Symbol   string `json:"symbol"`
		Interval string `json:"interval"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}
	
	if req.Strategy == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Strategy name required",
		})
	}
	
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	
	if req.Interval == "" {
		req.Interval = "15m"
	}
	
	// Check if already running
	if runningStrategies[req.Strategy] {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Strategy already running",
		})
	}
	
	// Mark as running
	runningStrategies[req.Strategy] = true
	strategyStopChans[req.Strategy] = make(chan bool)
	
	// Start strategy in background
	go runStrategyPaperTrading(req.Strategy, req.Symbol, req.Interval)
	
	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Strategy started",
		"strategy": req.Strategy,
		"symbol":   req.Symbol,
		"interval": req.Interval,
	})
}

// stopStrategyPaperTrading stops paper trading for a specific strategy
func stopStrategyPaperTrading(c *fiber.Ctx) error {
	var req struct {
		Strategy string `json:"strategy"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request",
		})
	}
	
	if req.Strategy == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Strategy name required",
		})
	}
	
	// Check if running
	if !runningStrategies[req.Strategy] {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Strategy not running",
		})
	}
	
	// Stop the strategy
	if stopChan, exists := strategyStopChans[req.Strategy]; exists {
		stopChan <- true
		close(stopChan)
		delete(strategyStopChans, req.Strategy)
	}
	
	runningStrategies[req.Strategy] = false
	
	return c.JSON(fiber.Map{
		"success":  true,
		"message":  "Strategy stopped",
		"strategy": req.Strategy,
	})
}

// runStrategyPaperTrading runs paper trading for a specific strategy
func runStrategyPaperTrading(strategy, symbol, interval string) {
	ticker := time.NewTicker(5 * time.Minute) // Check every 5 minutes
	defer ticker.Stop()
	
	stopChan := strategyStopChans[strategy]
	
	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			// Fetch candles
			candles, err := fetchBinanceData(symbol, interval, 4) // 4 days for indicators
			if err != nil {
				continue
			}
			
			// Generate signal for this strategy
			generator := &UnifiedSignalGenerator{}
			signal := generator.GenerateSignal(candles, strategy)
			
			if signal != nil && signal.Type != "NONE" {
				currentPrice := candles[len(candles)-1].Close
				paperTradingManager.AddTrade(signal, currentPrice)
			}
			
			// Update open trades
			currentPrice := candles[len(candles)-1].Close
			paperTradingManager.UpdateOpenTrades(currentPrice)
		}
	}
}
