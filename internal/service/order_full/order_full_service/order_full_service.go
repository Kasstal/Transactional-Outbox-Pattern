package order_full_service

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	outboxService "orders-center/internal/domain/outbox/service"
	paymentService "orders-center/internal/domain/payment/service"
	"orders-center/internal/service/order_full/entity"
)

type OrderFullService interface {
	GetOrderFull(ctx context.Context, id uuid.UUID) (entity.OrderFull, error)
}

type OrderFullSvc struct {
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	outboxService  outboxService.OutboxService
}

func NewOrderFullService(orderService orderService.OrderService,
	itemService itemService.OrderItemService,
	paymentService paymentService.PaymentService,
	historyService historyService.HistoryService,
	outboxService outboxService.OutboxService) OrderFullService {
	return &OrderFullSvc{
		orderService:   orderService,
		itemService:    itemService,
		paymentService: paymentService,
		historyService: historyService,
		outboxService:  outboxService,
	}
}

func (s *OrderFullSvc) GetOrderFull(ctx context.Context, id uuid.UUID) (entity.OrderFull, error) {
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
