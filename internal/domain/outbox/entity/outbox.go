package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type OutboxEvent struct {
	ID            uuid.UUID `json:"id"`
	AggregateType string    `json:"aggregate_type"` // Тип агрегата (например, "Order")
	AggregateID   uuid.UUID `json:"aggregate_id"`   // ID агрегата (например, ID заказа)
	EventType     string    `json:"event_type"`     // Тип события (например, "OrderCreated")
	Payload       []byte    `json:"payload"`        // Данные события (JSON)
	Status        string    `json:"status"`         // Статус события (pending, processed, failed)
	RetryCount    int32     `json:"retry_count"`    // Количество попыток
	CreatedAt     time.Time `json:"created_at"`
	ProcessedAt   time.Time `json:"processed_at,omitempty"` // Время обработки
	ErrorMessage  string    `json:"error_message,omitempty"`
}
