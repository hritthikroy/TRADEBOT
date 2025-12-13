package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthEndpoint(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/health", HealthHandler)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestReadinessEndpoint(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/ready", ReadinessHandler)

	req := httptest.NewRequest("GET", "/api/v1/ready", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 503 {
		t.Errorf("Expected status 200 or 503, got %d", resp.StatusCode)
	}
}

func TestLivenessEndpoint(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/live", LivenessHandler)

	req := httptest.NewRequest("GET", "/api/v1/live", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
