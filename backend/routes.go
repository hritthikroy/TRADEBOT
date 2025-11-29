package main

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Trading Bot API is running",
		})
	})

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
}
