package service

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/outbox/entity"
	"orders-center/internal/domain/outbox/repository"
	"time"
)

type OutboxService interface {
	AddEvent(ctx context.Context, event entity.OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error)
	MarkEventProcessed(ctx context.Context, eventID uuid.UUID) error
	UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error)
	AddNewEvent(ctx context.Context, event AddEventParams) error
}

type outboxService struct {
	repo repository.OutboxRepository
}

func NewOutboxService(repo repository.OutboxRepository) OutboxService {
	return &outboxService{repo: repo}
}

// Добавление события в Outbox
func (s *outboxService) AddEvent(ctx context.Context, event entity.OutboxEvent) error {
	event.Status = "pending"
	event.RetryCount = 0
	event.CreatedAt = time.Now()
	event.ID = uuid.Must(uuid.NewV4()) // Генерация уникального ID

	_, err := s.repo.CreateEvent(ctx, event)
	return err
}

type AddEventParams struct {
	AggregateType string          `json:"aggregate_type"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
}

func (s *outboxService) AddNewEvent(ctx context.Context, eventParams AddEventParams) error {
	event := entity.OutboxEvent{
		AggregateID:   eventParams.AggregateID,
		AggregateType: eventParams.AggregateType,
		EventType:     eventParams.EventType,
		Payload:       eventParams.Payload,
		Status:        "pending",
		RetryCount:    0,
	}

	_, err := s.repo.CreateEvent(ctx, event)
	return err
}

// Получение всех "pending" событий
func (s *outboxService) GetPendingEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	return s.repo.GetPendingEvents(ctx, limit)
}

// Обновление статуса события на "processed"
func (s *outboxService) MarkEventProcessed(ctx context.Context, eventID uuid.UUID) error {
	_, err := s.repo.UpdateEventStatus(ctx, eventID, "processed")
	return err
}

// Обновление статуса события на переданный статус
func (s *outboxService) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error) {
	return s.repo.UpdateEventStatus(ctx, eventID, status)
}
