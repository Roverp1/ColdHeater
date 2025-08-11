package database

import "time"

type Bot struct {
	Email        string     `db:"email"`
	IP           *string    `db:"ip"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	AgingEndDate *time.Time `db:"aging_end_date"`
}
