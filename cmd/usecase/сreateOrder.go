package usecase

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"orders-center/cmd/order_full/entity"
	transactional "orders-center/cmd/transactional"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	outboxService "orders-center/internal/domain/outbox/service"
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
	orderService orderService.OrderService,
	itemService itemService.OrderItemService,
	paymentService paymentService.PaymentService,
	historyService historyService.HistoryService,
	outboxService outboxService.OutboxService,
	txService *transactional.TransactionService,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderService:   orderService,
		itemService:    itemService,
		paymentService: paymentService,
		historyService: historyService,
		outboxService:  outboxService,
		txService:      txService,
	}
}

func (s *CreateOrderUseCase) Create(ctx context.Context, orderFull entity.OrderFull) error {

	err := s.txService.ExecTx(ctx, func(tx pgx.Tx) error {

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
