package handlers

import (
	"github.com/gofiber/fiber/v2"
	"tradebot/backend/internal/backtest"
)

// UnifiedBacktestRequest - Request for unified backtest
type UnifiedBacktestRequest struct {
	Symbol              string   `json:"symbol"`
	Interval            string   `json:"interval"`
	Days                int      `json:"days"`
	StartBalance        float64  `json:"startBalance"`
	Strategy            string   `json:"strategy"`
	
	// Optional advanced features
	EnableMonteCarlo    bool     `json:"enableMonteCarlo"`
	EnableStressTest    bool     `json:"enableStressTest"`
	EnableWalkForward   bool     `json:"enableWalkForward"`
	EnablePartialExits  bool     `json:"enablePartialExits"`
	EnableParallel      bool     `json:"enableParallel"`
	Strategies          []string `json:"strategies"`
	
	// Optional risk management
	RiskPercent         float64  `json:"riskPercent"`
	MaxDailyLoss        float64  `json:"maxDailyLoss"`
	MaxConsecutiveLoss  int      `json:"maxConsecutiveLoss"`
	MaxTradesPerDay     int      `json:"maxTradesPerDay"`
	
	// Optional filters
	TradingHoursOnly    bool     `json:"tradingHoursOnly"`
	MinVolatility       float64  `json:"minVolatility"`
	MaxVolatility       float64  `json:"maxVolatility"`
}

// HandleUnifiedBacktest - Single endpoint for all backtest needs
func HandleUnifiedBacktest(c *fiber.Ctx) error {
	var req UnifiedBacktestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request: " + err.Error(),
		})
	}
	
	// Validate required fields
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Interval == "" {
		req.Interval = "15m"
	}
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartBalance == 0 {
		req.StartBalance = 10000
	}
	if req.Strategy == "" {
		req.Strategy = "liquidity_hunter"
	}
	
	// Fetch candles
	candles, err := backtest.FetchBinanceData(req.Symbol, req.Interval, req.Days)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}
	
	if len(candles) < 100 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Insufficient data for backtest",
		})
	}
	
	// Build config
	config := backtest.UnifiedBacktestConfig{
		Symbol:              req.Symbol,
		Interval:            req.Interval,
		Days:                req.Days,
		StartBalance:        req.StartBalance,
		Strategy:            req.Strategy,
		EnableMonteCarlo:    req.EnableMonteCarlo,
		EnableStressTest:    req.EnableStressTest,
		UseWalkForward:      req.EnableWalkForward,
		EnablePartialExits:  req.EnablePartialExits,
		EnableParallel:      req.EnableParallel,
		Strategies:          req.Strategies,
		RiskPercent:         req.RiskPercent,
		MaxDailyLoss:        req.MaxDailyLoss,
		MaxConsecutiveLoss:  req.MaxConsecutiveLoss,
		MaxTradesPerDay:     req.MaxTradesPerDay,
		TradingHoursOnly:    req.TradingHoursOnly,
		MinVolatility:       req.MinVolatility,
		MaxVolatility:       req.MaxVolatility,
	}
	
	// Run unified backtest
	result, err := backtest.RunUnifiedBacktest(config, candles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Backtest failed: " + err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"result":  result,
	})
}
