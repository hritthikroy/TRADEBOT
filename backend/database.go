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
	
	// Supabase connection string
	connStr := fmt.Sprintf(
		"host=%s port=5432 user=postgres password=%s dbname=postgres sslmode=require",
		os.Getenv("SUPABASE_HOST"),
		os.Getenv("SUPABASE_PASSWORD"),
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Database connected successfully")
}
