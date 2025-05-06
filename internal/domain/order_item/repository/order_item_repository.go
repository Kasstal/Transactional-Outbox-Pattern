package repository

import (
	"context"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order_item/entity"
)

type orderItemRepository struct {
	q db.Queries
}

func newOrderItemRepository(q db.Queries) *orderItemRepository {
	return &orderItemRepository{q: q}
}

func (r *orderItemRepository) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (entity.OrderItem, error) {
	sqlArg := db.CreateOrderItemParams{
		ID:            arg.ID,
		ProductID:     arg.ProductID,
		ExternalID:    arg.ExternalID,
		Status:        arg.Status,
		BasePrice:     arg.BasePrice,
		Price:         arg.Price,
		EarnedBonuses: arg.EarnedBonuses,
		SpentBonuses:  arg.SpentBonuses,
		Gift:          arg.Gift,
		OwnerID:       arg.OwnerID,
		DeliveryID:    arg.DeliveryID,
		ShopAssistant: arg.ShopAssistant,
		Warehouse:     arg.Warehouse,
		OrderID:       arg.OrderID,
	}
	orderItem, err := r.q.CreateOrderItem(ctx, sqlArg)
	if err != nil {
		return entity.OrderItem{}, err
	}

	basePrice, err := orderItem.BasePrice.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	price, err := orderItem.Price.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	earnedBonuses, err := orderItem.EarnedBonuses.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	spentBonuses, err := orderItem.SpentBonuses.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}
	orderEntity := entity.OrderItem{
		ProductID:     orderItem.ProductID,
		ExternalID:    orderItem.ExternalID.String,
		Status:        orderItem.Status,
		BasePrice:     basePrice.Float64,
		Price:         price.Float64,
		EarnedBonuses: earnedBonuses.Float64,
		SpentBonuses:  spentBonuses.Float64,
		Gift:          orderItem.Gift.Bool,
		OwnerID:       orderItem.OwnerID.String,
		DeliveryID:    orderItem.DeliveryID.String,
		ShopAssistant: orderItem.ShopAssistant.String,
		Warehouse:     orderItem.Warehouse.String,
		OrderId:       orderItem.OrderID.Bytes,
	}

	return orderEntity, nil
}
func (r *orderItemRepository) GetOrderItem(ctx context.Context, id int32) (entity.OrderItem, error) {
	orderItem, err := r.q.GetOrderItem(ctx, id)

	if err != nil {
		return entity.OrderItem{}, err
	}

	basePrice, err := orderItem.BasePrice.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	price, err := orderItem.Price.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	earnedBonuses, err := orderItem.EarnedBonuses.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}

	spentBonuses, err := orderItem.SpentBonuses.Float64Value()
	if err != nil {
		return entity.OrderItem{}, err
	}
	orderEntity := entity.OrderItem{
		ProductID:     orderItem.ProductID,
		ExternalID:    orderItem.ExternalID.String,
		Status:        orderItem.Status,
		BasePrice:     basePrice.Float64,
		Price:         price.Float64,
		EarnedBonuses: earnedBonuses.Float64,
		SpentBonuses:  spentBonuses.Float64,
		Gift:          orderItem.Gift.Bool,
		OwnerID:       orderItem.OwnerID.String,
		DeliveryID:    orderItem.DeliveryID.String,
		ShopAssistant: orderItem.ShopAssistant.String,
		Warehouse:     orderItem.Warehouse.String,
		OrderId:       orderItem.OrderID.Bytes,
	}

	return orderEntity, nil

}
func (r *orderItemRepository) DeleteOrderItem(ctx context.Context, id int32) error {
	return r.q.DeleteOrderItem(ctx, id)
}
