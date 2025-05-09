package repository

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"orders-center/internal/domain/history/entity"
)

type CreateHistoryParams struct {
	Type     string          `json:"type"`
	TypeID   int32           `json:"type_id"`
	OldValue json.RawMessage `json:"old_value"`
	Value    json.RawMessage `json:"value"`
	UserID   string          `json:"user_id"`
	OrderID  uuid.UUID       `json:"order_id"`
}

type HistoryRepository interface {
	CreateHistory(ctx context.Context, arg CreateHistoryParams) (entity.History, error)
	GetHistory(ctx context.Context, id int32) (entity.History, error)
	DeleteHistory(ctx context.Context, id int32) error
	GetHistoriesByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.History, error)
}
