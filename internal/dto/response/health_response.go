package dto

import "time"

type HealthResponse struct {
	Status    string    `json:"status" example:"healthy"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Service   string    `json:"service" example:"api-account"`
}
