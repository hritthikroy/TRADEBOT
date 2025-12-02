package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Health check endpoints
	api.Get("/health", HealthHandler)
	api.Get("/ready", ReadinessHandler)
	api.Get("/live", LivenessHandler)

	// Signal routes
	signals := api.Group("/signals")
	signals.Post("/", CreateSignal)
	signals.Get("/", GetAllSignals)
	signals.Get("/pending", GetPendingSignals)
	signals.Get("/:id", GetSignalByID)
	signals.Put("/:id", UpdateSignal)
	signals.Put("/:id/live-price", UpdateLivePrice)
	signals.Delete("/:id", DeleteSignal)

	// Analytics routes
	analytics := api.Group("/analytics")
	analytics.Get("/", GetAnalytics)
	analytics.Get("/performance", GetPerformanceStats)
	analytics.Get("/by-killzone", GetStatsByKillZone)
	analytics.Get("/by-pattern", GetStatsByPattern)
	
	// Trade filtering routes
	filters := api.Group("/filters")
	filters.Get("/opportunities", GetBestTradeOpportunities) // Best trade opportunities
	filters.Get("/rules", GetTradeRules)                     // Trading rules
	filters.Get("/smart", GetSmartSignalFilter)              // Smart filter criteria
	
	// Backtest routes (API)
	backtest := api.Group("/backtest")
	backtest.Post("/run", HandleBacktestRunFiber)
	backtest.Post("/export", HandleBacktestExportFiber)
	backtest.Post("/comprehensive", HandleComprehensiveBacktest) // Test all strategies
	backtest.Post("/all-timeframes", HandleAllTimeframesBacktest) // Test all timeframes
	backtest.Post("/multi-timeframe", HandleMultiTimeframeBacktest) // Multi-TF top-down analysis
	
	// External Signal API routes (FREE)
	externalSignals := api.Group("/external-signals")
	externalSignals.Post("/get", HandleExternalSignals)      // Get external signals
	externalSignals.Post("/enhanced", HandleEnhancedSignal)  // Enhanced signal with external APIs
	externalSignals.Post("/compare", HandleCompareSignals)   // Compare multiple signals
	externalSignals.Get("/providers", HandleSignalProviders) // List free providers
	
	// Enhanced backtest with external signals
	backtest.Post("/enhanced", HandleEnhancedBacktest)    // Backtest with external signals

	// AI Enhancement routes
	ai := api.Group("/ai")
	ai.Get("/stats", GetAIStats)                          // Get AI statistics
	ai.Post("/toggle", ToggleAIFilter)                    // Enable/disable AI filter
	ai.Get("/test", TestAIConnection)                     // Test Grok AI connection
	ai.Post("/sentiment", AnalyzeSymbolSentiment)         // Analyze symbol sentiment

	// Template-based routes (Server-Side Rendered)
	app.Get("/", HandleIndexPage)
	app.Post("/backtest/run", HandleBacktestForm)
	app.Get("/signals/live", HandleLiveSignalsPage)

	// WebSocket route
	app.Use("/ws", WebSocketUpgrade)
	app.Get("/ws/signals", websocket.New(HandleWebSocket))
}
