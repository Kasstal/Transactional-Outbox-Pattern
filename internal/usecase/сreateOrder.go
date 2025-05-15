package usecase

import (
	"context"
	"log"
	historyService "orders-center/internal/domain/history/service"
	orderService "orders-center/internal/domain/order/service"
	itemService "orders-center/internal/domain/order_item/service"
	paymentService "orders-center/internal/domain/payment/service"
	"orders-center/internal/service/order_eno_1c"
	"orders-center/internal/service/order_full/entity"
	transactional "orders-center/internal/service/transactional"
)

type UseCase interface {
	Create(ctx context.Context, orderFull entity.OrderFull) error
}

type CreateOrderUseCase struct {
	enoService     *order_eno_1c.OrderEno1c
	orderService   orderService.OrderService
	itemService    itemService.OrderItemService
	paymentService paymentService.PaymentService
	historyService historyService.HistoryService
	txService      transactional.Transactional
}

func NewCreateOrderUseCase(
	enoService *order_eno_1c.OrderEno1c,
	orderService orderService.OrderService,
	itemService itemService.OrderItemService,
	paymentService paymentService.PaymentService,
	historyService historyService.HistoryService,
	txService transactional.Transactional,
) UseCase {
	return &CreateOrderUseCase{
		enoService:     enoService,
		orderService:   orderService,
		itemService:    itemService,
		paymentService: paymentService,
		historyService: historyService,
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

		err := s.enoService.CreateOutboxTask(ctx, orderFull.Order.ID)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	return err
}
