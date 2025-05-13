package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/history/entity"
)

type historyRepository struct {
	q db.Querier
}

func NewHistoryRepository(q db.Querier) HistoryRepository {
	return &historyRepository{q: q}
}

func (r *historyRepository) CreateHistory(ctx context.Context, arg CreateHistoryParams) (entity.History, error) {
	sqlArg := db.CreateHistoryParams{

		Type:     arg.Type,
		TypeID:   arg.TypeID,
		OldValue: arg.OldValue,
		Value:    arg.Value,
		UserID:   arg.UserID,
		OrderID:  pgtype.UUID{Bytes: arg.OrderID, Valid: true},
	}
	history, err := r.q.CreateHistory(ctx, sqlArg)
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
	history, err := r.q.GetHistory(ctx, id)
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
	return r.q.DeleteHistory(ctx, id)
}

func (r *historyRepository) GetHistoriesByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.History, error) {
	histories, err := r.q.GetHistoriesByOrderID(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
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
