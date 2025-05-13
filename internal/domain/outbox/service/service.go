package service

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"log"
	db "orders-center/db/sqlc"
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
	FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error)
	BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error)
	FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error)
	IncrementRetryCount(ctx context.Context, id uuid.UUID) error
}

type outboxService struct {
	repo repository.OutboxRepository
}

/*func NewOutboxService(repo repository.OutboxRepository) OutboxService {
	return &outboxService{repo: repo}
}*/

func NewOutboxService(q *db.Queries) OutboxService {

	repo := repository.NewOutboxRepository(q)
	return &outboxService{repo: repo}
}

func (s *outboxService) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	log.Println("entered outbox service")
	return s.repo.IncrementRetryCount(ctx, id)
}

func (s *outboxService) FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error) {
	return s.repo.FetchOnePendingForUpdateWithID(ctx, id)
}

func (s *outboxService) BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	return s.repo.BatchPendingTasks(ctx, limit)
}

func (s *outboxService) FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error) {
	return s.repo.FetchOnePendingForUpdate(ctx)
}

// Добавление события в Outbox
func (s *outboxService) AddEvent(ctx context.Context, event entity.OutboxEvent) error {
	event.Status = "pending"
	event.RetryCount = 0
	event.CreatedAt = time.Now()

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
