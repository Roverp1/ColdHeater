package database

import (
	"database/sql"
	"fmt"
	"time"

	"coldheater/internal/config"
)

func InsertBot(db *sql.DB, bot Bot, config *config.Config) error {
	agingEndDate := time.Now().AddDate(0, 0, config.Bot.AgingPeriod)

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

	row := db.QueryRow("SELECT email, ip, status, created_at, aging_end_date FROM bots WHERE email = $1", email)
	err := row.Scan(&bot.Email, &bot.IP, &bot.Status, &bot.CreatedAt, &bot.AgingEndDate)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan selected row %s into bot struct:\n%v", email, err)
	}

	return &bot, nil
}

func GetAllBots(db *sql.DB) ([]Bot, error) {
	var bots []Bot

	rows, err := db.Query("SELECT email, ip, status, created_at, aging_end_date FROM bots")
	if err != nil {
		return nil, fmt.Errorf("Failed to query all bots from database:\n%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bot Bot

		err := rows.Scan(&bot.Email, &bot.IP, &bot.Status, &bot.CreatedAt, &bot.AgingEndDate)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan selected row %s into bot struct:\n%v", bot.Email, err)
		}

		bots = append(bots, bot)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during rows iteration:\n%v", err)
	}

	return bots, err
}

func UpdateBotStatus(db *sql.DB, email, status string) error {
	result, err := db.Exec("UPDATE bots SET status = $1 WHERE email = $2", status, email)
	if err != nil {
		return fmt.Errorf("Failed to update bot's %s status:\n%v", email, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows after updating bot %s:\n%w", email, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No bot found with email %s: no rows affected", email)
	}

	return nil
}
