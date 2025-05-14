package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/inbox/entity"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/utils"
)

type inboxRepository struct {
	pool *pgxpool.Pool
}

func NewInboxRepository(pool *pgxpool.Pool) InboxRepository {
	return &inboxRepository{pool: pool}
}

func (r *inboxRepository) Create(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}
	sql_event_id, err := query.CreateInboxEvent(ctx, utils.ToUUID(event_id))
	if err != nil {
		return entity.InboxEvent{}, err
	}

	return entity.InboxEvent{EventID: sql_event_id.Bytes}, nil
}

func (r *inboxRepository) GetInboxEvent(ctx context.Context, event_id uuid.UUID) (entity.InboxEvent, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	sql_event_id, err := query.GetProcessedEvent(ctx, utils.ToUUID(event_id))
	if err != nil {
		return entity.InboxEvent{}, err
	}

	return entity.InboxEvent{EventID: sql_event_id.Bytes}, nil
}
