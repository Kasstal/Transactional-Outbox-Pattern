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
	FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error)
	BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error)
	FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error)
	IncrementRetryCount(ctx context.Context, id uuid.UUID) error
}
