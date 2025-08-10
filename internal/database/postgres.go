package database

import (
	"database/sql"
	"errors"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		return nil, errors.New("Mising DATABASE_URL env variable")
	}

	bd, err := sql.Open("postgres", DATABASE_URL)

	return bd, err
}
