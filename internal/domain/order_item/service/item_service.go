package service

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"orders-center/internal/domain/order_item/entity"
	"orders-center/internal/domain/order_item/repository"
	"orders-center/internal/utils"
)

type OrderItemService interface {
	GetByID(ctx context.Context, id int32) (entity.OrderItem, error)
	Create(ctx context.Context, item entity.OrderItem) (entity.OrderItem, error)
	GetOrderItemsByID(ctx context.Context, id uuid.UUID) ([]entity.OrderItem, error)
}

type orderItemService struct {
	repo repository.OrderItemRepository
}

func NewOrderItemService(repo repository.OrderItemRepository) OrderItemService {
	return &orderItemService{repo: repo}
}

// Получение OrderItem по ID
func (s *orderItemService) GetByID(ctx context.Context, id int32) (entity.OrderItem, error) {
	return s.repo.GetOrderItem(ctx, id)
}

// Создание нового OrderItem
func (s *orderItemService) Create(ctx context.Context, item entity.OrderItem) (entity.OrderItem, error) {
	// Преобразуем данные с учетом типов pgx
	arg := repository.CreateOrderItemParams{
		ID:            item.ID,
		ProductID:     item.ProductID,
		ExternalID:    utils.ToText(item.ExternalID), // Преобразуем ExternalID в pgtype.Text
		Status:        item.Status,
		BasePrice:     utils.ToNumeric(item.BasePrice),               // Преобразуем BasePrice в pgtype.Numeric
		Price:         utils.ToNumeric(item.Price),                   // Преобразуем Price в pgtype.Numeric
		EarnedBonuses: utils.ToNumeric(item.EarnedBonuses),           // Преобразуем EarnedBonuses в pgtype.Numeric
		SpentBonuses:  utils.ToNumeric(item.SpentBonuses),            // Преобразуем SpentBonuses в pgtype.Numeric
		Gift:          utils.ToBool(item.Gift),                       // Преобразуем Gift в pgtype.Bool
		OwnerID:       utils.ToText(item.OwnerID),                    // Преобразуем OwnerID в pgtype.Text
		DeliveryID:    utils.ToText(item.DeliveryID),                 // Преобразуем DeliveryID в pgtype.Text
		ShopAssistant: utils.ToText(item.ShopAssistant),              // Преобразуем ShopAssistant в pgtype.Text
		Warehouse:     utils.ToText(item.Warehouse),                  // Преобразуем Warehouse в pgtype.Text
		OrderID:       pgtype.UUID{Bytes: item.OrderId, Valid: true}, // Преобразуем OrderID в pgtype.UUID
	}

	// Сохраняем OrderItem в репозиторий
	return s.repo.CreateOrderItem(ctx, arg)
}

func (s *orderItemService) GetOrderItemsByID(ctx context.Context, id uuid.UUID) ([]entity.OrderItem, error) {
	return s.repo.GetOrderItemsByOrderID(ctx, id)
}
