package main

import (
	"github.com/gofiber/fiber/v2"
)

type TelegramBotStartRequest struct {
	Symbol   string `json:"symbol"`
	Strategy string `json:"strategy"`
	// Note: filterBuy and filterSell are no longer used here
	// They are controlled via user_settings API and read from database
}

// HandleStartTelegramBot starts the Telegram signal bot
func HandleStartTelegramBot(c *fiber.Ctx) error {
	var req TelegramBotStartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	// Validate inputs
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}
	if req.Strategy == "" {
		req.Strategy = "session_trader"
	}
	
	// Get current filter settings from database
	filterBuy, filterSell := GetCurrentFilterSettings()
	
	// Start the bot (pass current filter settings for display only)
	err := StartTelegramSignalBot(req.Symbol, req.Strategy, filterBuy, filterSell)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Telegram signal bot started successfully",
		"config": fiber.Map{
			"symbol":     req.Symbol,
			"strategy":   req.Strategy,
			"filterBuy":  filterBuy,  // Return current settings from database
			"filterSell": filterSell, // Return current settings from database
		},
	})
}

// HandleStopTelegramBot stops the Telegram signal bot
func HandleStopTelegramBot(c *fiber.Ctx) error {
	StopTelegramSignalBot()
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Telegram signal bot stopped successfully",
	})
}

// HandleGetTelegramBotStatus returns the bot status
func HandleGetTelegramBotStatus(c *fiber.Ctx) error {
	status := GetTelegramBotStatus()
	return c.JSON(status)
}
