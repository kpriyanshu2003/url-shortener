package database

import (
	"database/sql"
	"log"

	"github.com/kpriyanshu2003/url-shortener/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", config.GetEnv("PG_URL", "postgresql://postgres:admin@localhost/postgres?sslmode=disable"))
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL!")
}
