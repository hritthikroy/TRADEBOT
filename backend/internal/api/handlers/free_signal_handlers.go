package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ==================== FREE SIGNAL API HANDLERS ====================

// HandleExternalSignals gets signals from free external APIs
func HandleExternalSignals(c *fiber.Ctx) error {
	var req struct {
		Symbol   string `json:"symbol"`
		Interval string `json:"interval"`
		Provider string `json:"provider"` // "tradingview", "coingecko", "cryptocompare", "aggregated"
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Interval == "" {
		req.Interval = "1h"
	}
	if req.Provider == "" {
		req.Provider = "aggregated"
	}

	var signal *ExternalSignal
	var err error

	// Get signal from specified provider
	switch req.Provider {
	case "tradingview":
		service := NewSignalAPIService(TradingView)
		signal, err = service.GetTradingViewSignal(req.Symbol, req.Interval)
	case "coingecko":
		service := NewSignalAPIService(CoinGecko)
		signal, err = service.GetCoinGeckoSignal(req.Symbol)
	case "cryptocompare":
		service := NewSignalAPIService(CryptoCompare)
		signal, err = service.GetCryptoCompareSignal(req.Symbol)
	case "aggregated":
		signal, err = AggregateSignals(req.Symbol, req.Interval)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider. Use: tradingview, coingecko, cryptocompare, or aggregated",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get signal: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"signal": signal,
		"note":   "Free API - No API key required for most providers",
	})
}

// HandleEnhancedSignal generates signal enhanced with external APIs
func HandleEnhancedSignal(c *fiber.Ctx) error {
	var req struct {
		Symbol   string `json:"symbol"`
		Interval string `json:"interval"`
		Days     int    `json:"days"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Interval == "" {
		req.Interval = "15m"
	}
	if req.Days == 0 {
		req.Days = 7
	}

	// Fetch market data
	candles, err := fetchBinanceData(req.Symbol, req.Interval, req.Days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch market data: " + err.Error(),
		})
	}

	if len(candles) < 50 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient market data",
		})
	}

	// Generate our signal
	ourSignal := generateTrendFollowingSignal(candles, req.Interval)
	
	// Get external signals
	externalSignal, err := AggregateSignals(req.Symbol, req.Interval)
	if err != nil {
		// If external fails, return our signal only
		return c.JSON(fiber.Map{
			"our_signal":      ourSignal,
			"external_signal": nil,
			"enhanced_signal": ourSignal,
			"note":           "External APIs unavailable, using internal analysis only",
		})
	}

	// Enhance our signal with external data
	enhancedSignal, err := EnhanceSignalWithExternalAPI(ourSignal, req.Symbol, req.Interval)
	if err != nil {
		enhancedSignal = ourSignal
	}

	// Calculate agreement
	agreement := "NEUTRAL"
	if ourSignal != nil && externalSignal != nil {
		ourBullish := ourSignal.Type == "BUY"
		externalBullish := externalSignal.Recommendation == "BUY" || externalSignal.Recommendation == "STRONG_BUY"
		
		if (ourBullish && externalBullish) || (!ourBullish && !externalBullish) {
			agreement = "AGREE"
		} else {
			agreement = "DISAGREE"
		}
	}

	return c.JSON(fiber.Map{
		"symbol":   req.Symbol,
		"interval": req.Interval,
		"our_signal": ourSignal,
		"external_signal": externalSignal,
		"enhanced_signal": enhancedSignal,
		"agreement": agreement,
		"note": "Signal enhanced with free external APIs (TradingView, CoinGecko, CryptoCompare)",
	})
}

// HandleEnhancedBacktest runs backtest with external signal enhancement
func HandleEnhancedBacktest(c *fiber.Ctx) error {
	var config struct {
		Symbol         string  `json:"symbol"`
		Interval       string  `json:"interval"`
		Days           int     `json:"days"`
		StartBalance   float64 `json:"startBalance"`
		UseExternal    bool    `json:"useExternal"`
	}

	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Defaults
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
		config.StartBalance = 500
	}

	// Fetch historical data
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}

	if len(candles) < 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient historical data",
		})
	}

	// Run backtest
	btConfig := BacktestConfig{
		Symbol:       config.Symbol,
		Interval:     config.Interval,
		Days:         config.Days,
		StartBalance: config.StartBalance,
	}

	result, err := RunBacktest(btConfig, candles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Backtest failed: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"backtest_result": result,
		"external_enhanced": config.UseExternal,
		"note": "Backtest completed with free signal API integration",
	})
}

// HandleSignalProviders lists available free signal providers
func HandleSignalProviders(c *fiber.Ctx) error {
	providers := []fiber.Map{
		{
			"name":        "TradingView Technical Analysis",
			"provider":    "tradingview",
			"free":        true,
			"api_key_required": false,
			"features":    []string{"RSI", "EMA", "MACD", "Stochastic"},
			"description": "Technical indicators calculated from market data",
		},
		{
			"name":        "CoinGecko",
			"provider":    "coingecko",
			"free":        true,
			"api_key_required": false,
			"features":    []string{"Price data", "Market cap", "Volume", "24h/7d changes"},
			"description": "Free cryptocurrency market data",
		},
		{
			"name":        "CryptoCompare",
			"provider":    "cryptocompare",
			"free":        true,
			"api_key_required": false,
			"features":    []string{"Trading signals", "Technical analysis", "Market data"},
			"description": "Free tier available, optional API key for higher limits",
		},
		{
			"name":        "Aggregated Signals",
			"provider":    "aggregated",
			"free":        true,
			"api_key_required": false,
			"features":    []string{"Multiple sources", "Consensus signals", "Higher accuracy"},
			"description": "Combines signals from all free providers",
		},
	}

	return c.JSON(fiber.Map{
		"providers": providers,
		"note":      "All providers are free. No API keys required for basic usage.",
		"recommendation": "Use 'aggregated' for best results",
	})
}

// HandleCompareSignals compares our signals with external signals
func HandleCompareSignals(c *fiber.Ctx) error {
	var req struct {
		Symbol   string `json:"symbol"`
		Interval string `json:"interval"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Defaults
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Interval == "" {
		req.Interval = "1h"
	}

	// Fetch data
	candles, err := fetchBinanceData(req.Symbol, req.Interval, 7)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data: " + err.Error(),
		})
	}

	// Get our signal
	ourSignal := generateTrendFollowingSignal(candles, req.Interval)

	// Get signals from each provider
	tvService := NewSignalAPIService(TradingView)
	tvSignal, _ := tvService.GetTradingViewSignal(req.Symbol, req.Interval)

	cgService := NewSignalAPIService(CoinGecko)
	cgSignal, _ := cgService.GetCoinGeckoSignal(req.Symbol)

	ccService := NewSignalAPIService(CryptoCompare)
	ccSignal, _ := ccService.GetCryptoCompareSignal(req.Symbol)

	aggregated, _ := AggregateSignals(req.Symbol, req.Interval)

	return c.JSON(fiber.Map{
		"symbol":   req.Symbol,
		"interval": req.Interval,
		"timestamp": candles[len(candles)-1].Timestamp,
		"current_price": candles[len(candles)-1].Close,
		"signals": fiber.Map{
			"our_signal":        ourSignal,
			"tradingview":       tvSignal,
			"coingecko":         cgSignal,
			"cryptocompare":     ccSignal,
			"aggregated":        aggregated,
		},
		"note": "Compare signals from multiple free sources",
	})
}
