package domain

type Account struct {
	AccountID      int    `db:"account_id"`
	DocumentNumber string `db:"document_number"`
}
