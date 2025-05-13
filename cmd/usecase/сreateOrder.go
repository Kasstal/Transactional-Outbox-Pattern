package usecase

import (
	"context"
	"encoding/json"
	"log"
	"orders-center/cmd/order_full/entity"
	transactional "orders-center/cmd/transactional"
	db "orders-center/db/sqlc"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	outboxSvc "orders-center/internal/domain/outbox/service"
	paymentService "orders-center/internal/domain/payment/service"
)

type CreateOrderUseCase struct {
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	outboxService  outboxSvc.OutboxService
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

		orderService := orderService.NewOrderService(q)
		historyService := historyService.NewHistoryService(q)
		paymentService := paymentService.NewPaymentService(q)
		itemService := itemService.NewOrderItemService(q)
		outboxService := outboxSvc.NewOutboxService(q)
		if _, err := orderService.Create(ctx, orderFull.Order); err != nil {
			log.Println(err)
			return err
		}

		for _, item := range orderFull.Items {

			if _, err := itemService.Create(ctx, item); err != nil {
				log.Println(err)
				return err
			}
		}

		for _, payment := range orderFull.Payments {
			if _, err := paymentService.Create(ctx, payment); err != nil {
				log.Println(err)
				return err
			}
		}

		for _, history := range orderFull.History {
			if _, err := historyService.Create(ctx, history); err != nil {
				log.Println(err)
				return err
			}
		}

		eventData := map[string]interface{}{
			"order_id": orderFull.Order.ID,
		}
		payload, err := json.Marshal(eventData)
		if err != nil {
			log.Println(err)
			return err
		}
		// Добавляем событие в Outbox
		if err := outboxService.AddNewEvent(ctx, outboxSvc.AddEventParams{
			AggregateType: "OrderFull",        // Тип агрегата
			AggregateID:   orderFull.Order.ID, // ID заказа
			EventType:     "OrderCreated",     // Тип события
			Payload:       payload,            // Данные события
		}); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})
	log.Println(err)
	return err
}
