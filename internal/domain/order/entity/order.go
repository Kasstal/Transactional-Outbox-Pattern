package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	City        string    `json:"city"`
	Subdivision string    `json:"subdivision"`
	Price       float64   `json:"price"`
	Platform    string    `json:"platform"`
	GeneralID   uuid.UUID `json:"general_id"`
	OrderNumber string    `json:"order_number"`
	Executor    string    `json:"executor"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
