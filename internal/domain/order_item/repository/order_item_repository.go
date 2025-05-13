package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/order_item/entity"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/utils"
)

type orderItemRepository struct {
	pool *pgxpool.Pool
}

func NewOrderItemRepository(pool *pgxpool.Pool) OrderItemRepository {
	return &orderItemRepository{pool: pool}
}

func (r *orderItemRepository) GetOrderItemsByOrderID(ctx context.Context, id uuid.UUID) ([]entity.OrderItem, error) {
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	orderItems, err := query.GetOrderItemsByOrderID(ctx, utils.ToUUID(id))
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
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

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
	orderItem, err := query.CreateOrderItem(ctx, sqlArg)
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
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}

	orderItem, err := query.GetOrderItem(ctx, id)

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
	var query *db.Queries
	if tx, ok := transactional.TxFromContext(ctx); ok {
		query = db.New(tx)

	} else {
		query = db.New(r.pool)
	}
	return query.DeleteOrderItem(ctx, id)
}
