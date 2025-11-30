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
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	// Connect using the full connection string
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Printf("❌ Failed to ping database: %v", err)
		log.Printf("Connection string format: postgresql://user:password@host:port/database")
		log.Fatal("Database connection failed")
	}

	log.Println("✅ Database connected successfully")
}
