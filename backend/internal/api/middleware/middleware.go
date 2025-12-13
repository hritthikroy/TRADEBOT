package middleware

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// RequestValidationMiddleware validates common request parameters
func RequestValidationMiddleware(c *fiber.Ctx) error {
	// Add request ID for tracing
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = generateRequestID()
		c.Set("X-Request-ID", requestID)
	}

	c.Locals("requestID", requestID)
	return c.Next()
}

// MetricsMiddleware tracks request metrics
func MetricsMiddleware(c *fiber.Ctx) error {
	// Track request count, latency, etc.
	// This is a placeholder for future metrics implementation
	return c.Next()
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return generateUUID()
}

// generateUUID generates a simple UUID-like string
func generateUUID() string {
	b := make([]byte, 16)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// RecoverFromPanic recovers from panics in goroutines
func RecoverFromPanic(name string) {
	if r := recover(); r != nil {
		log.Printf("❌ Panic recovered in %s: %v", name, r)
	}
}
