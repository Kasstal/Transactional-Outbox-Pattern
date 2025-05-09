package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"orders-center/internal/domain/order/entity"
)

type CreateOrderParams struct {
	ID          uuid.UUID      `json:"id"`
	Type        string         `json:"type"`
	Status      string         `json:"status"`
	City        string         `json:"city"`
	Subdivision pgtype.Text    `json:"subdivision"`
	Price       pgtype.Numeric `json:"price"`
	Platform    string         `json:"platform"`
	GeneralID   pgtype.UUID    `json:"general_id"`
	OrderNumber string         `json:"order_number"`
	Executor    pgtype.Text    `json:"executor"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, arg CreateOrderParams) (entity.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (entity.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) error
}
