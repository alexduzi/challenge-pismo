package response

import "time"

type AccountResponse struct {
	AccountID      int64     `json:"account_id" example:"1"`
	DocumentNumber string    `json:"document_number" example:"12345678900"`
	FullName       string    `json:"full_name" example:"João da Silva"`
	Email          string    `json:"email" example:"joao.silva@example.com"`
	Phone          string    `json:"phone" example:"11987654321"`
	AccountType    string    `json:"account_type" example:"checking"`
	Balance        float64   `json:"balance" example:"1500.50"`
	CreatedAt      time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"2024-01-01T10:00:00Z"`
}
