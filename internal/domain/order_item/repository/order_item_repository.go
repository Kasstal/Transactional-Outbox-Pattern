package repository

import (
	"context"
	"github.com/gofrs/uuid"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order_item/entity"
	"orders-center/internal/utils"
)

type orderItemRepository struct {
	q *db.Queries
}

func NewOrderItemRepository(q *db.Queries) OrderItemRepository {
	return &orderItemRepository{q: q}
}

func (r *orderItemRepository) GetOrderItemsByOrderID(ctx context.Context, id uuid.UUID) ([]entity.OrderItem, error) {
	orderItems, err := r.q.GetOrderItemsByOrderID(ctx, utils.ToUUID(id))
	if err != nil {
		return nil, err
	}
	orderItemsEntity := make([]entity.OrderItem, len(orderItems))
	for _, orderItem := range orderItems {
		basePrice, err := orderItem.BasePrice.Float64Value()
		if err != nil {
			return []entity.OrderItem{}, err
		}

		price, err := orderItem.Price.Float64Value()
		if err != nil {
			return []entity.OrderItem{}, err
		}

		earnedBonuses, err := orderItem.EarnedBonuses.Float64Value()
		if err != nil {
			return []entity.OrderItem{}, err
		}

		spentBonuses, err := orderItem.SpentBonuses.Float64Value()
		if err != nil {
			return []entity.OrderItem{}, err
		}
		orderItemEntity := entity.OrderItem{
			ID:            orderItem.ID,
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
		orderItemsEntity = append(orderItemsEntity, orderItemEntity)

	}

	return orderItemsEntity, nil
}

func (r *orderItemRepository) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (entity.OrderItem, error) {
	sqlArg := db.CreateOrderItemParams{
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
		ID:            orderItem.ID,
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
