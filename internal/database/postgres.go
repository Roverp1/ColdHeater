package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		return nil, fmt.Errorf("Missing DATABASE_URL env variable")
	}

	db, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("Wasnt able to Open a database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to database: %v", err)
	}

	return db, nil
}
