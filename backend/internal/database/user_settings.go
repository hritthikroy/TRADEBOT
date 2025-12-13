package database

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

// UserSettings for API responses (camelCase for frontend)
type UserSettings struct {
	ID                 int      `json:"id"`
	FilterBuy          bool     `json:"filterBuy"`          // camelCase for frontend
	FilterSell         bool     `json:"filterSell"`         // camelCase for frontend
	SelectedStrategies []string `json:"selectedStrategies"` // camelCase for frontend
}

// UserSettingsDB for database operations (snake_case for Supabase)
type UserSettingsDB struct {
	ID                 int      `json:"id"`
	FilterBuy          bool     `json:"filter_buy"`          // snake_case for database
	FilterSell         bool     `json:"filter_sell"`         // snake_case for database
	SelectedStrategies []string `json:"selected_strategies"` // snake_case for database
}

// GetUserSettings retrieves user filter settings from Supabase
func GetUserSettings(c *fiber.Ctx) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		// Return defaults if Supabase not configured
		return c.JSON(UserSettings{
			ID:                 1,
			FilterBuy:          true,
			FilterSell:         true,
			SelectedStrategies: []string{"session_trader"},
		})
	}

	// Try to get settings from Supabase
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(UserSettings{
			ID:                 1,
			FilterBuy:          true,
			FilterSell:         true,
			SelectedStrategies: []string{"session_trader"},
		})
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Failed to get settings from Supabase: %v", err)
		return c.JSON(UserSettings{
			ID:                 1,
			FilterBuy:          true,
			FilterSell:         true,
			SelectedStrategies: []string{"session_trader"},
		})
	}
	defer resp.Body.Close()

	// Read response body for debugging
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	responseBodyStr := responseBody.String()
	
	log.Printf("🔍 GetUserSettings - Supabase response: %s", responseBodyStr)
	
	// Decode from database (snake_case)
	var settingsDB []UserSettingsDB
	if err := json.Unmarshal(responseBody.Bytes(), &settingsDB); err != nil {
		log.Printf("⚠️  GetUserSettings - Failed to decode: %v", err)
		log.Printf("⚠️  GetUserSettings - Response was: %s", responseBodyStr)
		return c.JSON(UserSettings{
			ID:                 1,
			FilterBuy:          true,
			FilterSell:         true,
			SelectedStrategies: []string{"session_trader"},
		})
	}

	if len(settingsDB) == 0 {
		// No settings exist, return defaults
		log.Printf("ℹ️  GetUserSettings - No settings found in database, returning defaults")
		return c.JSON(UserSettings{
			ID:                 1,
			FilterBuy:          true,
			FilterSell:         true,
			SelectedStrategies: []string{"session_trader"},
		})
	}

	// Convert from DB format to API format (snake_case to camelCase)
	apiSettings := UserSettings{
		ID:                 settingsDB[0].ID,
		FilterBuy:          settingsDB[0].FilterBuy,
		FilterSell:         settingsDB[0].FilterSell,
		SelectedStrategies: settingsDB[0].SelectedStrategies,
	}
	
	// Default to session_trader if no strategies selected
	if len(apiSettings.SelectedStrategies) == 0 {
		apiSettings.SelectedStrategies = []string{"session_trader"}
	}
	
	log.Printf("✅ GetUserSettings - Returning: filterBuy=%v, filterSell=%v, strategies=%v", apiSettings.FilterBuy, apiSettings.FilterSell, apiSettings.SelectedStrategies)
	return c.JSON(apiSettings)
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

	// Default to session_trader if no strategies provided
	if len(settings.SelectedStrategies) == 0 {
		settings.SelectedStrategies = []string{"session_trader"}
	}
	
	// Prepare data (only the fields we want to update)
	data := map[string]interface{}{
		"filter_buy":          settings.FilterBuy,
		"filter_sell":         settings.FilterSell,
		"selected_strategies": settings.SelectedStrategies,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal data",
		})
	}

	log.Printf("🔍 Updating settings in Supabase: %s", string(jsonData))

	// Use PATCH to update existing row
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Prefer", "return=representation")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error updating user settings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update settings",
		})
	}
	defer resp.Body.Close()

	// Read response to verify update
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	log.Printf("🔍 Supabase update response (status %d): %s", resp.StatusCode, responseBody.String())

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("❌ Supabase update failed with status %d", resp.StatusCode)
	}

	log.Printf("✅ User settings updated: filterBuy=%v, filterSell=%v, strategies=%v", settings.FilterBuy, settings.FilterSell, settings.SelectedStrategies)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Settings updated successfully",
		"settings": settings,
	})
}

// GetCurrentSettings returns current filter and strategy settings for internal use
func GetCurrentSettings() (bool, bool, []string) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Printf("⚠️  Supabase not configured, using default settings")
		return true, true, []string{"session_trader"}
	}

	// Try to get settings from Supabase
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("⚠️  Failed to create request, using defaults")
		return true, true, []string{"session_trader"}
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Failed to get settings from Supabase: %v, using defaults", err)
		return true, true, []string{"session_trader"}
	}
	defer resp.Body.Close()

	// Read response body
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	
	// Decode from database (snake_case)
	var settingsDB []UserSettingsDB
	if err := json.Unmarshal(responseBody.Bytes(), &settingsDB); err != nil {
		log.Printf("⚠️  Failed to decode settings: %v", err)
		return true, true, []string{"session_trader"}
	}

	if len(settingsDB) == 0 {
		log.Printf("ℹ️  No settings found in database, using defaults")
		return true, true, []string{"session_trader"}
	}

	strategies := settingsDB[0].SelectedStrategies
	if len(strategies) == 0 {
		strategies = []string{"session_trader"}
	}

	log.Printf("✅ Loaded settings: filterBuy=%v, filterSell=%v, strategies=%v", settingsDB[0].FilterBuy, settingsDB[0].FilterSell, strategies)
	return settingsDB[0].FilterBuy, settingsDB[0].FilterSell, strategies
}

// GetCurrentFilterSettings returns current filter settings for internal use (backward compatibility)
func GetCurrentFilterSettings() (bool, bool) {
	filterBuy, filterSell, _ := GetCurrentSettings()
	return filterBuy, filterSell
}

// Deprecated: Use GetCurrentSettings instead
func GetCurrentFilterSettingsOld() (bool, bool) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Printf("⚠️  Supabase not configured, using default filters (both enabled)")
		return true, true
	}

	// Try to get settings from Supabase (use headers for cache-busting, not query params)
	url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("⚠️  Failed to create request, using defaults")
		return true, true
	}

	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Expires", "0")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Failed to get settings from Supabase: %v, using defaults", err)
		return true, true
	}
	defer resp.Body.Close()

	// Read response body for debugging
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	responseBodyStr := responseBody.String()
	
	log.Printf("🔍 Supabase response for user_settings: %s", responseBodyStr)
	
	// Decode from database (snake_case)
	var settingsDB []UserSettingsDB
	if err := json.Unmarshal(responseBody.Bytes(), &settingsDB); err != nil {
		log.Printf("⚠️  Failed to decode settings: %v", err)
		log.Printf("⚠️  Response body was: %s", responseBodyStr)
		return true, true
	}

	if len(settingsDB) == 0 {
		log.Printf("ℹ️  No filter settings found in database, using defaults")
		return true, true
	}

	log.Printf("✅ Loaded filter settings from Supabase: filterBuy=%v, filterSell=%v", settingsDB[0].FilterBuy, settingsDB[0].FilterSell)
	return settingsDB[0].FilterBuy, settingsDB[0].FilterSell
}
