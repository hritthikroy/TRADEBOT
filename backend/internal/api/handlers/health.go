package handlers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

var startTime = time.Now()

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Version   string                 `json:"version"`
	Database  map[string]interface{} `json:"database"`
	System    SystemInfo             `json:"system"`
}

// SystemInfo represents system metrics
type SystemInfo struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"goroutines"`
	MemoryMB     uint64 `json:"memory_mb"`
	NumCPU       int    `json:"num_cpu"`
}

// HealthHandler returns detailed health information
func HealthHandler(c *fiber.Ctx) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Database:  GetDBStats(),
		System: SystemInfo{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			MemoryMB:     m.Alloc / 1024 / 1024,
			NumCPU:       runtime.NumCPU(),
		},
	}

	return c.JSON(health)
}

// ReadinessHandler checks if the service is ready to accept traffic
func ReadinessHandler(c *fiber.Ctx) error {
	// Check database connection
	if DB != nil {
		if err := DB.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status": "not_ready",
				"reason": "database_unavailable",
			})
		}
	}

	return c.JSON(fiber.Map{
		"status": "ready",
	})
}

// LivenessHandler checks if the service is alive
func LivenessHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "alive",
	})
}
