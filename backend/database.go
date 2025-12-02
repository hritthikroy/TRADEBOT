package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes database connection with retry logic and connection pooling
func InitDB() {
	host := os.Getenv("SUPABASE_HOST")
	password := os.Getenv("SUPABASE_PASSWORD")

	if host == "" || password == "" {
		log.Println("⚠️  Database credentials not configured")
		log.Println("⚠️  Running without database support")
		return
	}

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=postgres password=%s dbname=postgres sslmode=require",
		host, password,
	)

	// Retry logic
	maxRetries := 5
	var err error

	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("⚠️  Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Test connection
		err = DB.Ping()
		if err != nil {
			log.Printf("⚠️  Database ping attempt %d/%d failed: %v", i+1, maxRetries, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// Connection successful
		break
	}

	if err != nil {
		log.Printf("❌ Failed to connect to database after %d attempts: %v", maxRetries, err)
		log.Println("⚠️  Running without database support")
		DB = nil
		return
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("✅ Database connected successfully")
	
	// Start connection health monitor
	go monitorDBConnection()
}

// monitorDBConnection periodically checks database health
func monitorDBConnection() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if DB != nil {
			if err := DB.Ping(); err != nil {
				log.Printf("⚠️  Database health check failed: %v", err)
			}
		}
	}
}

// GetDBStats returns database connection pool statistics
func GetDBStats() map[string]interface{} {
	if DB == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	stats := DB.Stats()
	return map[string]interface{}{
		"status":          "connected",
		"open_connections": stats.OpenConnections,
		"in_use":          stats.InUse,
		"idle":            stats.Idle,
		"wait_count":      stats.WaitCount,
		"wait_duration":   stats.WaitDuration.String(),
		"max_idle_closed": stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}
}
