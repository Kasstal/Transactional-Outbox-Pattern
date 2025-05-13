package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/history/entity"
	transactional "orders-center/internal/service/transactional"
)

type historyRepository struct {
	pool *pgxpool.Pool
}

func NewHistoryRepository(pool *pgxpool.Pool) HistoryRepository {
	return &historyRepository{pool: pool}
}

func (r *historyRepository) CreateHistory(ctx context.Context, arg CreateHistoryParams) (entity.History, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	sqlArg := db.CreateHistoryParams{

		Type:     arg.Type,
		TypeID:   arg.TypeID,
		OldValue: arg.OldValue,
		Value:    arg.Value,
		UserID:   arg.UserID,
		OrderID:  pgtype.UUID{Bytes: arg.OrderID, Valid: true},
	}
	history, err := query.CreateHistory(ctx, sqlArg)
	if err != nil {
		return entity.History{}, err
	}
	return entity.History{
		Type:     history.Type,
		TypeId:   history.TypeID,
		OldValue: history.OldValue,
		Value:    history.Value,
		Date:     history.Date.Time,
		UserID:   history.UserID,
		OrderID:  history.OrderID.Bytes,
	}, nil
}
func (r *historyRepository) GetHistory(ctx context.Context, id int32) (entity.History, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}
	history, err := query.GetHistory(ctx, id)
	if err != nil {
		return entity.History{}, err
	}
	return entity.History{
		Type:     history.Type,
		TypeId:   history.TypeID,
		OldValue: history.OldValue,
		Value:    history.Value,
		Date:     history.Date.Time,
		UserID:   history.UserID,
		OrderID:  history.OrderID.Bytes,
	}, nil
}
func (r *historyRepository) DeleteHistory(ctx context.Context, id int32) error {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	return query.DeleteHistory(ctx, id)
}

func (r *historyRepository) GetHistoriesByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.History, error) {

	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	histories, err := query.GetHistoriesByOrderID(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if len(histories) == 0 {
		return []entity.History{}, fmt.Errorf("no history records with order_id = %d", orderID)
	}
	if err != nil {
		return []entity.History{}, err
	}
	historyEntities := make([]entity.History, len(histories))
	for _, history := range histories {
		historyEntity := entity.History{
			Type:     history.Type,
			TypeId:   history.TypeID,
			OldValue: history.OldValue,
			Value:    history.Value,
			Date:     history.Date.Time,
			UserID:   history.UserID,
			OrderID:  history.OrderID.Bytes,
		}
		historyEntities = append(historyEntities, historyEntity)
	}

	return historyEntities, nil
}
