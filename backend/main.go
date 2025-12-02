package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Trading Bot API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Serve static files from public directory
	app.Static("/", "../public")

	// Initialize database
	InitDB()

	// Start WebSocket hub
	go hub.Run()
	log.Println("✅ WebSocket hub started")

	// Start signal broadcaster (demo)
	StartSignalBroadcaster()
	log.Println("✅ Signal broadcaster started")

	// Routes
	SetupRoutes(app)

	// Start AI-enhanced signal generator (only if DB is available)
	if DB != nil {
		aiSignalGenerator := NewAIEnhancedSignalGenerator()
		aiSignalGenerator.Start()
		log.Println("✅ AI-Enhanced signal generator started")
		
		// Make it globally accessible for API endpoints
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("aiGenerator", aiSignalGenerator)
			return c.Next()
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
