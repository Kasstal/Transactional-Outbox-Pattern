package order_full_service

import (
	"context"
	"github.com/gofrs/uuid"
	"orders-center/cmd/order_full/entity"
	transactional "orders-center/cmd/transactional"
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
	txService      *transactional.TransactionService
}

func NewOrderFullService(orderService orderService.OrderService, itemService itemService.OrderItemService, paymentService paymentService.PaymentService, historyService historyService.HistoryService, outboxService outboxService.OutboxService, txService *transactional.TransactionService) *OrderFullService {
	return &OrderFullService{
		orderService:   orderService,
		itemService:    itemService,
		paymentService: paymentService,
		historyService: historyService,
		outboxService:  outboxService,
		txService:      txService,
	}
}

func (s *OrderFullService) GetOrderFull(ctx context.Context, id uuid.UUID) (entity.OrderFull, error) {
	// Собираем заказ
	order, err := s.orderService.GetByID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, err
	}

	// Собираем товары
	items, err := s.itemService.GetOrderItemsByID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, err
	}

	// Собираем платежи
	payments, err := s.paymentService.GetPaymentsByOrderID(ctx, id)
	if err != nil {
		return entity.OrderFull{}, err
	}

	// Собираем историю
	history, err := s.historyService.GetHistoriesByOrderId(ctx, id)
	if err != nil {
		return entity.OrderFull{}, err
	}

	// Возвращаем весь OrderFull
	return entity.OrderFull{
		Order:    order,
		Items:    items,
		Payments: payments,
		History:  history,
	}, nil
}
