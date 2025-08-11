package database

import (
	"database/sql"
	"fmt"
	"time"
)

func InsertBot(db *sql.DB, bot Bot) error {
	agingEndDate := time.Now().AddDate(0, 0, 14)

	var ipValue sql.NullString
	if bot.IP != nil {
		ipValue = sql.NullString{String: *bot.IP, Valid: true}
	} else {
		ipValue = sql.NullString{Valid: false}
	}

	_, err := db.Exec("INSERT INTO bots (email, status, ip, aging_end_date) VALUES ($1, $2, $3, $4)", bot.Email, bot.Status, ipValue, agingEndDate)
	if err != nil {
		return fmt.Errorf("Failed to insert bot %s:\n%v", bot.Email, err)
	}

	return nil
}

func GetBot(db *sql.DB, email string) (*Bot, error) {
	var bot Bot

	row := db.QueryRow("SELECT * FROM bots WHERE email = $1", email)
	err := row.Scan(&bot.Email, &bot.IP, &bot.Status, &bot.CreatedAt, &bot.AgingEndDate)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy seleted row into bot struct for %s:\n%v", email, err)
	}

	return &bot, nil
}
