package entity

import (
	"github.com/google/uuid"
	"time"
)

type History struct {
	Type     string    `json:"type"`
	TypeId   int32     `json:"type_id"`
	OldValue []byte    `json:"old_value"`
	Value    []byte    `json:"value"`
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	OrderID  uuid.UUID `json:"order_id"`
}
