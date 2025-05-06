package repository

import (
	"context"
	"fmt"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order/entity"
)

type orderRepository struct {
	q db.Queries
}

func NewOrderRepository(q db.Queries) OrderRepository {
	return &orderRepository{q: q}
}
func (r *orderRepository) CreateOrder(ctx context.Context, arg CreateOrderParams) (entity.Order, error) {
	sqlArg := db.CreateOrderParams{
		ID:          arg.ID,
		Type:        arg.Type,
		Status:      arg.Status,
		City:        arg.City,
		Subdivision: arg.Subdivision,
		Price:       arg.Price,
		Platform:    arg.Platform,
		GeneralID:   arg.GeneralID,
		OrderNumber: arg.OrderNumber,
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
		ID:          fmt.Sprint(order.ID),
		Type:        order.Type,
		Status:      order.Status,
		City:        order.City,
		Subdivision: order.Subdivision.String,
		Price:       price.Float64,
		Platform:    order.Platform,
		GeneralID:   order.GeneralID.Bytes,
		OrderNumber: order.OrderNumber,
		Executor:    order.Executor.String,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
	}
	return orderEntity, nil
}
func (r *orderRepository) GetOrder(ctx context.Context, id int32) (entity.Order, error) {
	order, err := r.q.GetOrder(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	price, err := order.Price.Float64Value()
	if err != nil {
		return entity.Order{}, err
	}
	orderEntity := entity.Order{
		ID:          fmt.Sprint(order.ID),
		Type:        order.Type,
		Status:      order.Status,
		City:        order.City,
		Subdivision: order.Subdivision.String,
		Price:       price.Float64,
		Platform:    order.Platform,
		GeneralID:   order.GeneralID.Bytes,
		OrderNumber: order.OrderNumber,
		Executor:    order.Executor.String,
		CreatedAt:   order.CreatedAt.Time,
		UpdatedAt:   order.UpdatedAt.Time,
	}
	return orderEntity, nil
}
func (r *orderRepository) DeleteOrder(ctx context.Context, id int32) error {
	return r.q.DeleteOrder(ctx, id)
}
