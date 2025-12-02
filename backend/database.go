package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	
	// Get DATABASE_URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	
	if databaseURL == "" {
		log.Println("⚠️  DATABASE_URL not set - running in backtest-only mode")
		return
	}

	// Connect using the full connection string
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		log.Println("⚠️  Running in backtest-only mode")
		return
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Printf("❌ Failed to ping database: %v", err)
		log.Printf("Connection string format: postgresql://user:password@host:port/database")
		log.Println("⚠️  Running in backtest-only mode")
		DB = nil
		return
	}

	log.Println("✅ Database connected successfully")
}
