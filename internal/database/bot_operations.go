package database

import (
	"database/sql"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"coldheater/internal/config"
)

func InsertBot(db *sql.DB, bot *Bot) error {
	var agingEndDate time.Time
	if bot.AgingEndDate != nil {
		agingEndDate = *bot.AgingEndDate
	} else {
		agingEndDate = time.Now().AddDate(0, 0, config.Global.Bot.AgingPeriod)
	}

	columns := []string{"email", "aging_end_date", "status", "created_at", "last_used", "first_name", "last_name", "username", "password"}
	values := []any{bot.Email, agingEndDate, bot.Status, bot.CreatedAt, bot.LastUsed, bot.FirstName, bot.LastName, bot.Username, bot.Password}
	var placeholders []string
	placeHolderIndex := len(columns)

	for i := 1; i <= placeHolderIndex; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}

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

func GetVerificationAccount(db *sql.DB) (*VerificationAccount, error) {
	var verificationAccs []VerificationAccount
	var randomAcc *VerificationAccount

	rows, err := db.Query("SELECT email, password, created_at, last_used, usage_count, is_active FROM verification_accounts WHERE is_active = true")
	if err != nil {
		return nil, fmt.Errorf("Failed to select all verification accounts from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var acc VerificationAccount

		err := rows.Scan(&acc.Email, &acc.Password, &acc.CreatedAt, &acc.LastUsed, &acc.UsageCount, &acc.IsActive)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan select row into VerificationAccount struct: %w", err)
		}

		verificationAccs = append(verificationAccs, acc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during VerificationAccount's rows iteration: %w", err)
	}

	if len(verificationAccs) == 0 {
		return nil, fmt.Errorf("No active verification accounts found")
	}

	randomAcc = &verificationAccs[rand.IntN(len(verificationAccs))]
	return randomAcc, nil
}

func IncrementVerificationAccUsage(db *sql.DB, email string) error {
	result, err := db.Exec("UPDATE verification_accounts SET usage_count = usage_count + 1, last_used = NOW() WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("Failed to update verification acc's %s usage_count: %w", email, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to get affected rows after updating verification account %s:\n%w", email, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No verification account found with email %s: no rows affected", email)
	}

	return nil
}
