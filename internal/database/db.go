package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kpriyanshu2003/url-shortener/internal/utils"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		utils.GetEnv("PG_HOST", "localhost"),
		utils.GetEnv("PG_PORT", "5432"),
		utils.GetEnv("PG_USER", "postgres"),
		utils.GetEnv("PG_PASSWORD", "admin"),
		utils.GetEnv("PG_DB_NAME", "url_shortener"),
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
