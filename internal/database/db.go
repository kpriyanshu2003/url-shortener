package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kpriyanshu2003/url-shortener/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("PG_HOST", "localhost"),
		config.GetEnv("PG_PORT", "5432"),
		config.GetEnv("PG_USER", "postgres"),
		config.GetEnv("PG_PASSWORD", "admin"),
		config.GetEnv("PG_DB_NAME", "postgres"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL!")
}
