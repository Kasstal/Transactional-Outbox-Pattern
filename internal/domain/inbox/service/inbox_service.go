package service

import (
	"context"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/inbox/entity"
	"orders-center/internal/domain/inbox/repository"
)

type InboxService interface {
	Create(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error)
	GetInboxEvent(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error)
}

type inboxService struct {
	repo repository.InboxRepository
}

func NewInboxService(repo repository.InboxRepository) InboxService {
	return &inboxService{repo: repo}
}
func (i inboxService) Create(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error) {
	return i.repo.Create(ctx, event_id)
}

func (i inboxService) GetInboxEvent(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error) {
	return i.repo.GetInboxEvent(ctx, event_id)
}
