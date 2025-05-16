package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/outbox/entity"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/utils"
)

type outboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) OutboxRepository {
	return &outboxRepository{
		pool: pool,
	}
}

func (r *outboxRepository) MarkStatusFailed(ctx context.Context, id uuid.UUID, errorMsg string) error {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	params := db.MarkEventErrorParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		ErrorMessage: pgtype.Text{String: errorMsg, Valid: true},
	}

	err := query.MarkEventError(ctx, params)
	if err != nil {
		return err
	}
	return nil

}

func (r *outboxRepository) GetAllInProgressEvents(ctx context.Context) ([]entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	events, err := query.GetAllInProgressOutboxEvents(ctx)
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
			ErrorMessage:  event.ErrorMessage.String,
		})
	}

	return result, nil
}

func (r *outboxRepository) BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	events, err := query.BatchPendingTasks(ctx, int32(limit))

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
			ErrorMessage:  event.ErrorMessage.String,
		})
	}
	return result, nil
}

func (r *outboxRepository) IncrementRetryCount(ctx context.Context, id uuid.UUID, errMsg string) (int32, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	log.Println("entered outboxrepo")
	params := db.IncrementRetryCountParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		ErrorMessage: pgtype.Text{String: errMsg, Valid: true},
	}

	retryCount, err := query.IncrementRetryCount(ctx, params)
	if err == nil {
		log.Println("increased retry count")
	}
	return retryCount.Int32, err
}

func (r *outboxRepository) FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	event, err := query.FetchOnePendingForUpdateWithID(ctx, utils.ToUUID(id))
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
		ErrorMessage:  event.ErrorMessage.String,
	}, nil
}

func (r *outboxRepository) FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	event, err := query.FetchOnePendingForUpdate(ctx)
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
		ErrorMessage:  event.ErrorMessage.String,
	}, nil

}

// Создание нового события в Outbox
func (r *outboxRepository) CreateEvent(ctx context.Context, event entity.OutboxEvent) (entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	sqlArg := db.CreateOutboxEventParams{
		AggregateType: event.AggregateType,
		AggregateID:   utils.ToUUID(event.AggregateID),
		EventType:     event.EventType,
		Payload:       event.Payload,
		Status:        event.Status,
		RetryCount:    pgtype.Int4{int32(event.RetryCount), true},
	}

	sqlEvent, err := query.CreateOutboxEvent(ctx, sqlArg)
	if err != nil {
		return entity.OutboxEvent{}, err
	}

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
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	events, err := query.GetPendingOutboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	var result []entity.OutboxEvent
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
			ErrorMessage:  event.ErrorMessage.String,
		})
	}

	return result, nil
}

// Обновление статуса события в Outbox
func (r *outboxRepository) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	sqlArg := db.UpdateOutboxEventStatusParams{
		ID:     utils.ToUUID(eventID),
		Status: status,
	}

	sqlEvent, err := query.UpdateOutboxEventStatus(ctx, sqlArg)
	if err != nil {
		return entity.OutboxEvent{}, err
	}

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
		ErrorMessage:  sqlEvent.ErrorMessage.String,
	}, nil
}

// Удаление события из Outbox
func (r *outboxRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	err := query.DeleteOutboxEvent(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
	return err
}
