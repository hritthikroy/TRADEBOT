package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HandleGetRecentSignals retrieves recent signals from Supabase
func HandleGetRecentSignals(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	
	signals, err := GetRecentSignals(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"count":   len(signals),
		"signals": signals,
	})
}

// HandleGetSignalPerformance retrieves performance metrics
func HandleGetSignalPerformance(c *fiber.Ctx) error {
	performance, err := GetSignalPerformance()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success":     true,
		"performance": performance,
	})
}

// HandleUpdateSignalStatus updates a signal's status
func HandleUpdateSignalStatus(c *fiber.Ctx) error {
	type UpdateRequest struct {
		SignalID          string  `json:"signal_id"`
		Status            string  `json:"status"`
		CurrentPrice      float64 `json:"current_price"`
		ProfitLoss        float64 `json:"profit_loss"`
		ProfitLossPercent float64 `json:"profit_loss_percent"`
	}
	
	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	
	err := UpdateSignalStatus(req.SignalID, req.Status, req.CurrentPrice, req.ProfitLoss, req.ProfitLossPercent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Signal status updated successfully",
	})
}
