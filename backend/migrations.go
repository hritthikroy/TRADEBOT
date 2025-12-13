package main

import (
	"log"
)

// RunMigrations runs database migrations
func RunMigrations() {
	if DB == nil {
		// No direct database connection - using Supabase REST API instead
		// Migrations are handled via Supabase SQL Editor
		return
	}

	log.Println("🔄 Running database migrations...")

	// Create user_settings table if it doesn't exist
	createUserSettingsTable()

	log.Println("✅ Database migrations complete")
}

// createUserSettingsTable creates the user_settings table
func createUserSettingsTable() {
	query := `
		CREATE TABLE IF NOT EXISTS user_settings (
			id INTEGER PRIMARY KEY DEFAULT 1,
			filter_buy BOOLEAN DEFAULT true,
			filter_sell BOOLEAN DEFAULT true,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT single_row CHECK (id = 1)
		);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("⚠️  Failed to create user_settings table: %v", err)
		return
	}

	log.Println("✅ user_settings table ready")

	// Insert default settings if not exists
	insertQuery := `
		INSERT INTO user_settings (id, filter_buy, filter_sell)
		SELECT 1, true, true
		WHERE NOT EXISTS (SELECT 1 FROM user_settings WHERE id = 1);
	`

	_, err = DB.Exec(insertQuery)
	if err != nil {
		log.Printf("⚠️  Failed to insert default settings: %v", err)
		return
	}

	log.Println("✅ Default filter settings initialized")
}
