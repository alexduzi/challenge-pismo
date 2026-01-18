package domain

import "time"

type Account struct {
	AccountID      int64     `db:"account_id"`
	DocumentNumber string    `db:"document_number"`
	FullName       string    `db:"full_name"`
	Email          string    `db:"email"`
	Phone          string    `db:"phone"`
	AccountType    string    `db:"account_type"`
	Balance        float64   `db:"balance"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
