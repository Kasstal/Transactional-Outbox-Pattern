package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/outbox/entity"
	"orders-center/internal/utils"
)

type outboxRepository struct {
	q *db.Queries
}

func NewOutboxRepository(q *db.Queries) OutboxRepository {
	return &outboxRepository{
		q: q,
	}
}

func (r *outboxRepository) BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	events, err := r.q.BatchPendingTasks(ctx, int32(limit))

	if err != nil {
		return nil, err
	}

	result := make([]entity.OutboxEvent, len(events))
	for _, event := range events {
		result = append(result, entity.OutboxEvent{
			ID:            event.ID.Bytes,
			AggregateType: event.AggregateType,
			AggregateID:   event.AggregateID.Bytes,
			EventType:     event.EventType,
			Payload:       event.Payload,
			Status:        event.Status,
			RetryCount:    event.RetryCount.Int32,
			CreatedAt:     event.CreatedAt.Time,
			ProcessedAt:   event.ProcessedAt.Time,
		})
	}
	return result, nil
}

func (r *outboxRepository) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	log.Println("entered outboxrepo")
	err := r.q.IncrementRetryCount(ctx, utils.ToUUID(id))
	if err == nil {
		log.Println("increased retry count")
	}
	return err
}

func (r *outboxRepository) FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error) {
	event, err := r.q.FetchOnePendingForUpdateWithID(ctx, utils.ToUUID(id))
	if err != nil {
		return entity.OutboxEvent{}, err
	}
	return entity.OutboxEvent{
		ID:            event.ID.Bytes,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID.Bytes,
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		RetryCount:    event.RetryCount.Int32,
		CreatedAt:     event.CreatedAt.Time,
		ProcessedAt:   event.ProcessedAt.Time,
	}, nil
}

func (r *outboxRepository) FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error) {
	event, err := r.q.FetchOnePendingForUpdate(ctx)
	if err != nil {
		return entity.OutboxEvent{}, err
	}
	return entity.OutboxEvent{
		ID:            event.ID.Bytes,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID.Bytes,
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		RetryCount:    event.RetryCount.Int32,
		CreatedAt:     event.CreatedAt.Time,
		ProcessedAt:   event.ProcessedAt.Time,
	}, nil

}

// Создание нового события в Outbox
func (r *outboxRepository) CreateEvent(ctx context.Context, event entity.OutboxEvent) (entity.OutboxEvent, error) {
	sqlArg := db.CreateOutboxEventParams{
		AggregateType: event.AggregateType,
		AggregateID:   utils.ToUUID(event.AggregateID),
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		RetryCount:    pgtype.Int4{int32(event.RetryCount), true},
	}

	// Вызов SQL-запроса
	sqlEvent, err := r.q.CreateOutboxEvent(ctx, sqlArg)
	if err != nil {
		return entity.OutboxEvent{}, err
	}

	// Преобразуем данные обратно в структуру домена
	return entity.OutboxEvent{
		ID:            sqlEvent.ID.Bytes,
		AggregateType: sqlEvent.AggregateType,
		AggregateID:   sqlEvent.AggregateID.Bytes,
		EventType:     sqlEvent.EventType,
		Payload:       sqlEvent.Payload,
		Status:        sqlEvent.Status,
		RetryCount:    sqlEvent.RetryCount.Int32,
		CreatedAt:     sqlEvent.CreatedAt.Time,
		ProcessedAt:   sqlEvent.ProcessedAt.Time,
	}, nil
}

// Получение "pending" событий
func (r *outboxRepository) GetPendingEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	events, err := r.q.GetPendingOutboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	var result []entity.OutboxEvent
	for _, event := range events {
		// Преобразуем данные обратно в структуру домена
		result = append(result, entity.OutboxEvent{
			ID:            event.ID.Bytes,
			AggregateType: event.AggregateType,
			AggregateID:   event.AggregateID.Bytes,
			EventType:     event.EventType,
			Payload:       event.Payload,
			Status:        event.Status,
			RetryCount:    event.RetryCount.Int32,
			CreatedAt:     event.CreatedAt.Time,
			ProcessedAt:   event.ProcessedAt.Time,
		})
	}

	return result, nil
}

// Обновление статуса события в Outbox
func (r *outboxRepository) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error) {
	sqlArg := db.UpdateOutboxEventStatusParams{
		ID:     utils.ToUUID(eventID),
		Status: status,
	}

	// Выполняем SQL-запрос
	sqlEvent, err := r.q.UpdateOutboxEventStatus(ctx, sqlArg)
	if err != nil {
		return entity.OutboxEvent{}, err
	}

	// Преобразуем данные обратно в структуру домена
	return entity.OutboxEvent{
		ID:            sqlEvent.ID.Bytes,
		AggregateType: sqlEvent.AggregateType,
		AggregateID:   sqlEvent.AggregateID.Bytes,
		EventType:     sqlEvent.EventType,
		Payload:       sqlEvent.Payload,
		Status:        sqlEvent.Status,
		RetryCount:    sqlEvent.RetryCount.Int32,
		CreatedAt:     sqlEvent.CreatedAt.Time,
		ProcessedAt:   sqlEvent.ProcessedAt.Time,
	}, nil
}

// Удаление события из Outbox
func (r *outboxRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	err := r.q.DeleteOutboxEvent(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	return err
}
