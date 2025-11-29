package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	
	// Get credentials from environment
	host := os.Getenv("SUPABASE_HOST")
	password := os.Getenv("SUPABASE_PASSWORD")
	
	// Build connection string using key=value format (safer for special characters)
	connStr := fmt.Sprintf(
		"host=%s port=5432 user=postgres password=%s dbname=postgres sslmode=require",
		host,
		password,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection with timeout
	err = DB.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		log.Printf("Host: %s", host)
		log.Fatal("Database connection failed")
	}

	log.Println("✅ Database connected successfully")
}
