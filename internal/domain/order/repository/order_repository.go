package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order/entity"
	"orders-center/internal/utils"
)

type orderRepository struct {
	q db.Queries
}

func NewOrderRepository(q db.Queries) OrderRepository {
	return &orderRepository{q: q}
}

func (r *orderRepository) CreateOrder(ctx context.Context, arg CreateOrderParams) (entity.Order, error) {
	sqlArg := db.CreateOrderParams{
		ID:          pgtype.UUID{Bytes: arg.ID, Valid: true},
		Type:        arg.Type,
		Status:      arg.Status,
		City:        arg.City,
		Subdivision: arg.Subdivision,
		Price:       arg.Price,
		Platform:    arg.Platform,
		GeneralID:   arg.GeneralID,
		OrderNumber: utils.ToText(arg.OrderNumber),
		Executor:    arg.Executor,
	}
	order, err := r.q.CreateOrder(ctx, sqlArg)
	if err != nil {
		return entity.Order{}, err
	}
	price, err := order.Price.Float64Value()
	if err != nil {
		return entity.Order{}, err
	}
	orderEntity := entity.Order{
		ID:          order.ID.Bytes,
		Type:        order.Type,
		Status:      order.Status,
		City:        order.City,
		Subdivision: order.Subdivision.String,
		Price:       price.Float64,
		Platform:    order.Platform,
		GeneralID:   order.GeneralID.Bytes,
		OrderNumber: order.OrderNumber.String,
		Executor:    order.Executor.String,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
	}
	return orderEntity, nil
}
func (r *orderRepository) GetOrder(ctx context.Context, id uuid.UUID) (entity.Order, error) {
	order, err := r.q.GetOrder(ctx, pgtype.UUID{Bytes: id})
	if err != nil {
		return entity.Order{}, err
	}
	price, err := order.Price.Float64Value()
	if err != nil {
		return entity.Order{}, err
	}
	orderEntity := entity.Order{
		ID:          order.ID.Bytes,
		Type:        order.Type,
		Status:      order.Status,
		City:        order.City,
		Subdivision: order.Subdivision.String,
		Price:       price.Float64,
		Platform:    order.Platform,
		GeneralID:   order.GeneralID.Bytes,
		OrderNumber: order.OrderNumber.String,
		Executor:    order.Executor.String,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
	}
	return orderEntity, nil
}
func (r *orderRepository) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteOrder(ctx, pgtype.UUID{Bytes: id})
}
