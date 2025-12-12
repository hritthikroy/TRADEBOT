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
	backtest.Post("/test-all-strategies", HandleTestAllStrategies) // Test 10 advanced strategies
	backtest.Post("/optimize-parameters", HandleOptimizeParameters) // Optimize single strategy parameters
	backtest.Post("/optimize-all", HandleOptimizeAllStrategies) // Optimize all strategies
	backtest.Post("/world-class-optimize", HandleWorldClassOptimization) // World-class optimization
	backtest.Post("/quick-optimize", HandleQuickOptimization) // Quick single strategy optimization
	backtest.Post("/live-signal", HandleLiveSignalFiber) // Get live trading signal
	backtest.Post("/world-class", HandleWorldClassBacktest) // World-class backtest with advanced metrics
	backtest.Post("/compare", HandleQuickCompare)           // Compare standard vs world-class backtest
	backtest.Post("/ai-optimize", HandleAIOptimization)     // AI-powered parameter optimization
	backtest.Get("/ai-analyze", HandleAIMarketAnalysis)     // AI market analysis
	backtest.Get("/ai-recommend", HandleAIStrategyRecommendation) // AI strategy recommendation
	backtest.Get("/ai-compare", HandleAICompareStrategies)  // AI-powered strategy comparison
	
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
	
	// Telegram Bot routes
	telegram := api.Group("/telegram")
	telegram.Post("/start", HandleStartTelegramBot)       // Start 24/7 signal bot
	telegram.Post("/stop", HandleStopTelegramBot)         // Stop signal bot
	telegram.Get("/status", HandleGetTelegramBotStatus)   // Get bot status
	
	// Signal Storage routes
	signalStorage := api.Group("/signal-storage")
	signalStorage.Get("/recent", HandleGetRecentSignals)        // Get recent signals
	signalStorage.Get("/performance", HandleGetSignalPerformance) // Get performance metrics
	signalStorage.Post("/update", HandleUpdateSignalStatus)      // Update signal status
	
	// User Settings routes
	settings := api.Group("/settings")
	settings.Get("/", GetUserSettings)           // Get filter settings
	settings.Post("/", UpdateUserSettings)       // Update filter settings

	// Paper Trading routes
	paperTrading := api.Group("/paper-trading")
	paperTrading.Get("/stats", getPaperTradingStats)           // Get paper trading statistics
	paperTrading.Get("/trades", getAllPaperTrades)             // Get all paper trades
	paperTrading.Post("/trade", addPaperTrade)                 // Add new paper trade
	paperTrading.Post("/update", updatePaperTrades)            // Update open trades
	paperTrading.Post("/reset", resetPaperTrading)             // Reset paper trading
	paperTrading.Post("/start-auto", startAutoPaperTrading)    // Start auto paper trading
	paperTrading.Post("/stop-auto", stopAutoPaperTrading)      // Stop auto paper trading
	paperTrading.Post("/add-test-trade", addTestTrade)         // Add test trade for demo

	// Template-based routes (Server-Side Rendered)
	app.Get("/", HandleIndexPage)
	app.Post("/backtest/run", HandleBacktestForm)
	app.Get("/signals/live", HandleLiveSignalsPage)
	app.Get("/paper-trading", HandlePaperTradingPage)

	// WebSocket route
	app.Use("/ws", WebSocketUpgrade)
	app.Get("/ws/signals", websocket.New(HandleWebSocket))
}
