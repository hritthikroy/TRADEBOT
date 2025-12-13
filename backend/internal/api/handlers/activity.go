package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// HandleGetActivities returns recent activities
func HandleGetActivities(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 100)
	activityType := c.Query("type", "")
	level := c.Query("level", "")

	logger := GetActivityLogger()
	var activities []Activity

	if activityType != "" {
		activities = logger.GetActivitiesByType(ActivityType(activityType), limit)
	} else if level != "" {
		activities = logger.GetActivitiesByLevel(ActivityLevel(level), limit)
	} else {
		activities = logger.GetActivities(limit)
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"activities": activities,
		"count":      len(activities),
	})
}

// HandleGetActivityStats returns activity statistics
func HandleGetActivityStats(c *fiber.Ctx) error {
	logger := GetActivityLogger()
	stats := logger.GetStats()

	return c.JSON(fiber.Map{
		"success": true,
		"stats":   stats,
	})
}

// HandleClearActivities clears all activities
func HandleClearActivities(c *fiber.Ctx) error {
	logger := GetActivityLogger()
	logger.Clear()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Activity log cleared",
	})
}

// HandleActivityWebSocket handles WebSocket connections for real-time activity updates
func HandleActivityWebSocket(c *websocket.Conn) {
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
	logger := GetActivityLogger()
	
	// Subscribe to activity updates
	activityChan := logger.Subscribe(clientID)
	defer logger.Unsubscribe(clientID)

	// Send initial activities
	activities := logger.GetActivities(50)
	if err := c.WriteJSON(fiber.Map{
		"type":       "initial",
		"activities": activities,
	}); err != nil {
		return
	}

	// Listen for new activities and send to client
	done := make(chan struct{})
	
	// Goroutine to read from WebSocket (to detect disconnection)
	go func() {
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				close(done)
				return
			}
		}
	}()

	// Send activities to client
	for {
		select {
		case activity := <-activityChan:
			if err := c.WriteJSON(fiber.Map{
				"type":     "activity",
				"activity": activity,
			}); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}

// HandleStartActivityRecording starts recording all activities
func HandleStartActivityRecording(c *fiber.Ctx) error {
	LogSystemSuccess("Activity recording started", fiber.Map{
		"timestamp": time.Now(),
		"user":      c.IP(),
	})

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Activity recording started",
	})
}

// HandleStopActivityRecording stops recording activities
func HandleStopActivityRecording(c *fiber.Ctx) error {
	logger := GetActivityLogger()
	stats := logger.GetStats()

	LogSystemSuccess("Activity recording stopped", stats)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Activity recording stopped",
		"stats":   stats,
	})
}

// HandleExportActivities exports activities as JSON
func HandleExportActivities(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 1000)
	logger := GetActivityLogger()
	activities := logger.GetActivities(limit)

	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=activities_%s.json", time.Now().Format("2006-01-02_15-04-05")))

	return c.JSON(fiber.Map{
		"exported_at": time.Now(),
		"count":       len(activities),
		"activities":  activities,
	})
}
