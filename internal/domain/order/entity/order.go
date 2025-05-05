package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Order struct {
	ID          string
	Type        string
	Status      string
	City        string
	Subdivision string
	Price       float64
	Platform    string
	GeneralID   uuid.UUID
	OrderNumber string
	Executor    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
