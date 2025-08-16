package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"coldheater/internal/config"
)

func InsertBot(db *sql.DB, bot Bot, config *config.Config) error {
	var agingEndDate time.Time
	if bot.AgingEndDate != nil {
		agingEndDate = *bot.AgingEndDate
	} else {
		agingEndDate = time.Now().AddDate(0, 0, config.Bot.AgingPeriod)
	}

	columns := []string{"email", "aging_end_date"}
	placeholders := []string{"$1", "$2"}
	values := []any{bot.Email, agingEndDate}
	var placeHolderIndex uint8 = 3

	insertBotHelper := func(value *string, insert string) {
		if value != nil {
			values = append(values, value)
			placeholders = append(placeholders, fmt.Sprintf("$%d", placeHolderIndex))
			columns = append(columns, insert)
			placeHolderIndex++
		}
	}

	insertBotHelper(bot.Status, "status")
	insertBotHelper(bot.FirstName, "first_name")
	insertBotHelper(bot.LastName, "last_name")
	insertBotHelper(bot.Username, "username")
	insertBotHelper(bot.Password, "password")

	sqlQuery := fmt.Sprintf("INSERT INTO bots (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err := db.Exec(sqlQuery, values...)
	if err != nil {
		return fmt.Errorf("Failed to insert bot %s:\n%v", bot.Email, err)
	}

	return nil
}

func GetBot(db *sql.DB, email string) (*Bot, error) {
	var bot Bot

	row := db.QueryRow("SELECT email, status, created_at, last_used, aging_end_date, first_name, last_name, username, password FROM bots WHERE email = $1", email)
	err := row.Scan(&bot.Email, &bot.Status, &bot.CreatedAt, &bot.LastUsed, &bot.AgingEndDate, &bot.FirstName, &bot.LastName, &bot.Username, &bot.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan selected row %s into bot struct:\n%v", email, err)
	}

	return &bot, nil
}

func GetAllBots(db *sql.DB) ([]Bot, error) {
	var bots []Bot

	rows, err := db.Query("SELECT email, status, created_at, last_used, aging_end_date, first_name, last_name, username, password FROM bots")
	if err != nil {
		return nil, fmt.Errorf("Failed to query all bots from database:\n%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bot Bot

		err := rows.Scan(&bot.Email, &bot.Status, &bot.CreatedAt, &bot.LastUsed, &bot.AgingEndDate, &bot.FirstName, &bot.LastName, &bot.Username, &bot.Password)
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
