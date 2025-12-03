package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserSettings struct {
	ID         int  `json:"id"`
	FilterBuy  bool `json:"filterBuy"`
	FilterSell bool `json:"filterSell"`
}

// GetUserSettings retrieves user filter settings from Supabase
func GetUserSettings(c *fiber.Ctx) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		// Return defaults if Supabase not configured
		return c.JSON(UserSettings{
			ID:         1,
			FilterBuy:  true,
			FilterSell: true,
		})
	}

	// Try to get settings from Supabase
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(UserSettings{ID: 1, FilterBuy: true, FilterSell: true})
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Failed to get settings from Supabase: %v", err)
		return c.JSON(UserSettings{ID: 1, FilterBuy: true, FilterSell: true})
	}
	defer resp.Body.Close()

	var settings []UserSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return c.JSON(UserSettings{ID: 1, FilterBuy: true, FilterSell: true})
	}

	if len(settings) == 0 {
		// No settings exist, return defaults
		return c.JSON(UserSettings{ID: 1, FilterBuy: true, FilterSell: true})
	}

	return c.JSON(settings[0])
}

// UpdateUserSettings updates user filter settings in Supabase
func UpdateUserSettings(c *fiber.Ctx) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	var settings UserSettings
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if supabaseURL == "" || supabaseKey == "" {
		log.Printf("⚠️  Supabase not configured, settings not persisted")
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Settings updated (not persisted - Supabase not configured)",
			"settings": settings,
		})
	}

	// Prepare data
	data := map[string]interface{}{
		"id":          1,
		"filter_buy":  settings.FilterBuy,
		"filter_sell": settings.FilterSell,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal data",
		})
	}

	// Upsert to Supabase
	url := fmt.Sprintf("%s/rest/v1/user_settings", supabaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error updating user settings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update settings",
		})
	}
	defer resp.Body.Close()

	log.Printf("✅ User settings updated: filterBuy=%v, filterSell=%v", settings.FilterBuy, settings.FilterSell)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Settings updated successfully",
		"settings": settings,
	})
}

// GetCurrentFilterSettings returns current filter settings for internal use
func GetCurrentFilterSettings() (bool, bool) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Printf("⚠️  Supabase not configured, using default filters (both enabled)")
		return true, true
	}

	// Try to get settings from Supabase with cache-busting
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1&_=%d", supabaseURL, time.Now().UnixNano())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("⚠️  Failed to create request, using defaults")
		return true, true
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Failed to get settings from Supabase: %v, using defaults", err)
		return true, true
	}
	defer resp.Body.Close()

	var settings []UserSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		log.Printf("⚠️  Failed to decode settings, using defaults")
		return true, true
	}

	if len(settings) == 0 {
		log.Printf("ℹ️  No filter settings found, using defaults")
		return true, true
	}

	log.Printf("✅ Loaded filter settings from Supabase: filterBuy=%v, filterSell=%v", settings[0].FilterBuy, settings[0].FilterSell)
	return settings[0].FilterBuy, settings[0].FilterSell
}
