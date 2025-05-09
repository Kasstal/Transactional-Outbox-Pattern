package repository

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"orders-center/internal/domain/order_item/entity"
)

type CreateOrderItemParams struct {
	ID            int32          `json:"id"`
	ProductID     string         `json:"product_id"`
	ExternalID    pgtype.Text    `json:"external_id"`
	Status        string         `json:"status"`
	BasePrice     pgtype.Numeric `json:"base_price"`
	Price         pgtype.Numeric `json:"price"`
	EarnedBonuses pgtype.Numeric `json:"earned_bonuses"`
	SpentBonuses  pgtype.Numeric `json:"spent_bonuses"`
	Gift          pgtype.Bool    `json:"gift"`
	OwnerID       pgtype.Text    `json:"owner_id"`
	DeliveryID    pgtype.Text    `json:"delivery_id"`
	ShopAssistant pgtype.Text    `json:"shop_assistant"`
	Warehouse     pgtype.Text    `json:"warehouse"`
	OrderID       pgtype.UUID    `json:"order_id"`
}
type OrderItemRepository interface {
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (entity.OrderItem, error)
	GetOrderItem(ctx context.Context, id int32) (entity.OrderItem, error)
	DeleteOrderItem(ctx context.Context, id int32) error
	GetOrderItemsByOrderID(ctx context.Context, id uuid.UUID) ([]entity.OrderItem, error)
}
