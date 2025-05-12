package usecase

import (
	"context"
	"encoding/json"
	"orders-center/cmd/order_full/entity"
	transactional "orders-center/cmd/transactional"
	db "orders-center/db/sqlc"
	historyRepo "orders-center/internal/domain/history/repository"
	historyService "orders-center/internal/domain/history/service"
	orderRepo "orders-center/internal/domain/order/repository"
	orderService "orders-center/internal/domain/order/service"
	itemRepo "orders-center/internal/domain/order_item/repository"
	itemService "orders-center/internal/domain/order_item/service"
	outboxRepo "orders-center/internal/domain/outbox/repository"
	outboxService "orders-center/internal/domain/outbox/service"
	paymentRepo "orders-center/internal/domain/payment/repository"
	paymentService "orders-center/internal/domain/payment/service"
)

type CreateOrderUseCase struct {
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	outboxService  outboxService.OutboxService
	txService      *transactional.TransactionService
}

func NewCreateOrderUseCase(
	/*orderService orderService.OrderService,
	  itemService itemService.OrderItemService,
	  paymentService paymentService.PaymentService,
	  historyService historyService.HistoryService,
	  outboxService outboxService.OutboxService,*/
	txService *transactional.TransactionService,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		/*orderService:   orderService,
		itemService:    itemService,
		paymentService: paymentService,
		historyService: historyService,
		outboxService:  outboxService,*/
		txService: txService,
	}
}

func (s *CreateOrderUseCase) Create(ctx context.Context, orderFull entity.OrderFull) error {

	err := s.txService.ExecTx(ctx, func(q *db.Queries) error {

		query := q
		orderRepo := orderRepo.NewOrderRepository(query)
		itemRepo := itemRepo.NewOrderItemRepository(query)
		paymentRepo := paymentRepo.NewPaymentRepository(query)
		historyRepo := historyRepo.NewHistoryRepository(query)
		outboxRepo := outboxRepo.NewOutboxRepository(query)
		s.orderService = orderService.NewOrderService(orderRepo)
		s.itemService = itemService.NewOrderItemService(itemRepo)
		s.paymentService = paymentService.NewPaymentService(paymentRepo)
		s.historyService = historyService.NewHistoryService(historyRepo)
		s.outboxService = outboxService.NewOutboxService(outboxRepo)

		if _, err := s.orderService.Create(ctx, orderFull.Order); err != nil {
			return err
		}

		for _, item := range orderFull.Items {
			if _, err := s.itemService.Create(ctx, item); err != nil {
				return err
			}
		}

		for _, payment := range orderFull.Payments {
			if _, err := s.paymentService.Create(ctx, payment); err != nil {
				return err
			}
		}

		for _, history := range orderFull.History {
			if _, err := s.historyService.Create(ctx, history); err != nil {
				return err
			}
		}

		eventData := map[string]interface{}{
			"order_id": orderFull.Order.ID,
		}
		payload, err := json.Marshal(eventData)
		if err != nil {
			return err
		}
		// Добавляем событие в Outbox
		if err := s.outboxService.AddNewEvent(ctx, outboxService.AddEventParams{
			AggregateType: "OrderFull",        // Тип агрегата
			AggregateID:   orderFull.Order.ID, // ID заказа
			EventType:     "OrderCreated",     // Тип события
			Payload:       payload,            // Данные события
		}); err != nil {
			return err
		}

		return nil
	})

	return err
}
