package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// Validate required environment variables
	validateEnv()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Trading Bot API v1.0",
		ErrorHandler: customErrorHandler,
	})

	// Panic recovery middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Structured logging middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
	}))

	// Rate limiting middleware
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
		},
	}))

	// CORS middleware - allow all origins for development
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: false,
	}))

	// Serve static files from public directory
	app.Static("/", "../public")

	// Initialize database
	InitDB()

	// Run database migrations
	RunMigrations()

	// Start WebSocket hub
	go hub.Run()
	log.Println("✅ WebSocket hub started")

	// Start signal broadcaster (demo)
	StartSignalBroadcaster()
	log.Println("✅ Signal broadcaster started")
	
	// Initialize Telegram bot
	InitTelegramBot()
	
	// Auto-start Telegram bot if configured
	autoStart := os.Getenv("TELEGRAM_AUTO_START")
	if autoStart == "true" {
		symbol := os.Getenv("TELEGRAM_SYMBOL")
		strategy := os.Getenv("TELEGRAM_STRATEGY")
		filterBuy := os.Getenv("TELEGRAM_FILTER_BUY") == "true"
		filterSell := os.Getenv("TELEGRAM_FILTER_SELL") == "true"
		
		if symbol == "" {
			symbol = "BTCUSDT"
		}
		if strategy == "" {
			strategy = "session_trader"
		}
		
		err := StartTelegramSignalBot(symbol, strategy, filterBuy, filterSell)
		if err != nil {
			log.Printf("⚠️  Failed to auto-start Telegram bot: %v", err)
		} else {
			log.Printf("✅ Telegram bot auto-started for %s with %s strategy", symbol, strategy)
		}
	}

	// Start AI-enhanced signal generator (only if DB is available)
	// Must be initialized BEFORE routes so middleware is available
	if DB != nil {
		aiSignalGenerator := NewAIEnhancedSignalGenerator()
		aiSignalGenerator.Start()
		log.Println("✅ AI-Enhanced signal generator started")
		
		// Make it globally accessible for API endpoints
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("aiGenerator", aiSignalGenerator)
			return c.Next()
		})
	} else {
		// If DB is not available, provide a nil-safe middleware
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("aiGenerator", nil)
			return c.Next()
		})
	}

	// Routes (must be after middleware setup)
	SetupRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("⚠️  PORT not set, defaulting to 8080")
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📊 Dashboard: http://localhost:%s", port)
	log.Printf("🏥 Health: http://localhost:%s/api/v1/health", port)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// validateEnv checks required environment variables
func validateEnv() {
	// Check if Supabase REST API is configured (preferred method)
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	if supabaseURL != "" && supabaseKey != "" {
		log.Println("✅ Supabase REST API configured")
		return
	}
	
	// Check if direct PostgreSQL connection is configured (alternative method)
	dbHost := os.Getenv("SUPABASE_HOST")
	dbPassword := os.Getenv("SUPABASE_PASSWORD")
	
	if dbHost != "" && dbPassword != "" {
		log.Println("✅ PostgreSQL direct connection configured")
		return
	}
	
	// Neither method configured
	log.Println("⚠️  No database configured - some features will use defaults")
	log.Println("ℹ️  To enable persistence, configure SUPABASE_URL and SUPABASE_KEY in .env")
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("❌ Error: %v (Path: %s)", err, c.Path())

	return c.Status(code).JSON(fiber.Map{
		"error":   message,
		"code":    code,
		"path":    c.Path(),
		"message": fmt.Sprintf("%v", err),
	})
}
