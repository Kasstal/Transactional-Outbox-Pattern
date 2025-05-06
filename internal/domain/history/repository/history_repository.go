package repository

import (
	"context"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/history/entity"
)

type historyRepository struct {
	q db.Queries
}

func NewHistoryRepository(q db.Queries) HistoryRepository {
	return &historyRepository{q: q}
}

func (r *historyRepository) CreateHistory(ctx context.Context, arg CreateHistoryParams) (entity.History, error) {
	sqlArg := db.CreateHistoryParams{
		Type:     arg.Type,
		TypeID:   arg.TypeID,
		OldValue: arg.OldValue,
		Value:    arg.Value,
		UserID:   arg.UserID,
		OrderID:  arg.OrderID,
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
		OrderID:  history.OrderID,
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
		OrderID:  history.OrderID,
	}, nil
}
func (r *historyRepository) DeleteHistory(ctx context.Context, id int32) error {
	return r.q.DeleteHistory(ctx, id)
}
