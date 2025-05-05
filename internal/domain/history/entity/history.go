package entity

import "time"

type History struct {
	Type     string    `json:"type"`
	TypeId   int       `json:"type_id"`
	OldValue []byte    `json:"old_value"`
	Value    []byte    `json:"value"`
	Date     time.Time `json:"date"`
	UserID   string    `json:"user_id"`
	OrderID  string    `json:"order_id"`
}
