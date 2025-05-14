package usecase

import (
	"context"
	"encoding/json"
	"log"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	outboxSvc "orders-center/internal/domain/outbox/service"
	paymentService "orders-center/internal/domain/payment/service"
	"orders-center/internal/service/order_full/entity"
	transactional "orders-center/internal/service/transactional"
)

type UseCase interface {
	Create(ctx context.Context, orderFull entity.OrderFull) error
}

type CreateOrderUseCase struct {
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	outboxService  outboxSvc.OutboxService
	txService      transactional.Transactional
}

func NewCreateOrderUseCase(
	orderService orderService.OrderService,
	itemService itemService.OrderItemService,
	paymentService paymentService.PaymentService,
	historyService historyService.HistoryService,
	outboxService outboxSvc.OutboxService,
	txService transactional.Transactional,
) UseCase {
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

	err := s.txService.ExecTx(ctx, func(ctx context.Context) error {

		if _, err := s.orderService.Create(ctx, orderFull.Order); err != nil {
			log.Println(err)
			return err
		}

		for _, item := range orderFull.Items {

			if _, err := s.itemService.Create(ctx, item); err != nil {
				log.Println(err)
				return err
			}
		}

		for _, payment := range orderFull.Payments {
			if _, err := s.paymentService.Create(ctx, payment); err != nil {
				log.Println(err)
				return err
			}
		}

		for _, history := range orderFull.History {
			if _, err := s.historyService.Create(ctx, history); err != nil {
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
		// Add event в Outbox
		if err := s.outboxService.AddNewEvent(ctx, outboxSvc.AddEventParams{
			AggregateType: "OrderFull",
			AggregateID:   orderFull.Order.ID,
			EventType:     "OrderCreated",
			Payload:       payload,
		}); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})
	log.Println(err)
	return err
}
