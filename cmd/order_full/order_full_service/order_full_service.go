package order_full_service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"orders-center/cmd/order_full/entity"
	db "orders-center/db/sqlc"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	outboxService "orders-center/internal/domain/outbox/service"
	paymentService "orders-center/internal/domain/payment/service"
)

type OrderFullService struct {
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	outboxService  outboxService.OutboxService
}

func NewOrderFullService(q *db.Queries) *OrderFullService {
	return &OrderFullService{
		orderService:   orderService.NewOrderService(q),
		itemService:    itemService.NewOrderItemService(q),
		paymentService: paymentService.NewPaymentService(q),
		historyService: historyService.NewHistoryService(q),
		outboxService:  outboxService.NewOutboxService(q),
	}
}

func (s *OrderFullService) GetOrderFull(ctx context.Context, id uuid.UUID) (entity.OrderFull, error) {
	// Собираем заказ
	order, err := s.orderService.GetByID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, fmt.Errorf("from order: ", err)
	}
	// Собираем товары
	items, err := s.itemService.GetOrderItemsByID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, fmt.Errorf("from items: ", err)
	}

	// Собираем платежи
	payments, err := s.paymentService.GetPaymentsByOrderID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, fmt.Errorf("from payments: ", err)
	}

	// Собираем историю
	history, err := s.historyService.GetHistoriesByOrderId(ctx, id)
	if err != nil {
		return entity.OrderFull{}, fmt.Errorf("from history: ", err)
	}

	// Возвращаем весь OrderFull
	return entity.OrderFull{
		Order:    order,
		Items:    items,
		Payments: payments,
		History:  history,
	}, nil
}
