package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/inbox/entity"
)

type InboxRepository interface {
	Create(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error)
	GetInboxEvent(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error)
}
