package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/outbox/entity"
)

// Интерфейс репозитория для работы с таблицей outbox_events
type OutboxRepository interface {
	CreateEvent(ctx context.Context, event entity.OutboxEvent) (entity.OutboxEvent, error)
	GetPendingEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error)
	UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
}
