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
	
	// Supabase connection string with proper format
	// Format: postgres://postgres:[PASSWORD]@[HOST]:5432/postgres
	connStr := fmt.Sprintf(
		"postgres://postgres:%s@%s:5432/postgres?sslmode=require",
		os.Getenv("SUPABASE_PASSWORD"),
		os.Getenv("SUPABASE_HOST"),
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Database connected successfully")
}
