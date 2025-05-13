package service

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order/entity"
	"orders-center/internal/domain/order/repository"
	"orders-center/internal/utils"
)

type OrderService interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error)
	Create(ctx context.Context, order entity.Order) (entity.Order, error)
	// Update(ctx context.Context, order entity.Order) error // TODO: позже
}

type orderService struct {
	repo repository.OrderRepository
}

/*
func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}*/

func NewOrderService(q *db.Queries) OrderService {
	repo := repository.NewOrderRepository(q)
	return &orderService{repo: repo}
}

func (s *orderService) GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *orderService) Create(ctx context.Context, order entity.Order) (entity.Order, error) {
	arg := repository.CreateOrderParams{
		ID:          order.ID,
		Type:        order.Type,
		Status:      order.Status,
		City:        order.City,
		Subdivision: pgtype.Text{String: order.Subdivision, Valid: order.Subdivision != ""},
		Price:       utils.ToNumeric(order.Price),
		Platform:    order.Platform,
		GeneralID:   pgtype.UUID{Bytes: order.GeneralID, Valid: true},
		OrderNumber: order.OrderNumber,
		Executor:    pgtype.Text{String: order.Executor, Valid: order.Executor != ""},
	}
	return s.repo.CreateOrder(ctx, arg)
}
