package database

import "time"

type Bot struct {
	Email        string     `db:"email"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	LastUsed     time.Time  `db:"last_used"`
	AgingEndDate *time.Time `db:"aging_end_date"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Username     string     `db:"username"`
	Password     string     `db:"password"`
	BirthDate    time.Time  `db:"birth_date"`
	Gender       string     `db:"gender"`
}

type VerificationAccount struct {
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	CreatedAt  time.Time  `db:"created_at"`
	LastUsed   *time.Time `db:"last_used"`
	UsageCount int        `db:"usage_count"`
	IsActive   bool       `db:"is_active"`
}

type Sender struct {
	Id 				int	
	Email			string
	Daily_quota 	string
	Warmup_end_date time.Time	
	Campaign_id 	int
}
