package order_eno_1c

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"orders-center/cmd/cron"
	orderFullService "orders-center/cmd/order_full/order_full_service"
	transactional "orders-center/cmd/transactional"
	db "orders-center/db/sqlc"
	outboxRepo "orders-center/internal/domain/outbox/repository"
	outboxService "orders-center/internal/domain/outbox/service"
	"orders-center/internal/utils"
	"time"
)

type OrderEno1c struct {
	cron                 cron.Scheduler
	outboxService        outboxService.OutboxService
	transactionalService *transactional.TransactionService
	orderFullService     *orderFullService.OrderFullService
}

func NewOrderEno1c(cron cron.Scheduler, outbox outboxService.OutboxService, transactionalService *transactional.TransactionService, orderFullService *orderFullService.OrderFullService) *OrderEno1c {
	return &OrderEno1c{
		cron:                 cron,
		outboxService:        outbox,
		transactionalService: transactionalService,
		orderFullService:     orderFullService,
	}
}

func (o *OrderEno1c) Run(ctx context.Context) {
	o.cron.AddFunc()
	o.cron.Start(ctx, o.processTask)
}

func (o *OrderEno1c) ProcessTask(ctx context.Context) error {
	//timeout := 5 * time.Second
	//taskCtx, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()

	err := o.transactionalService.ExecTx(ctx, func(q *db.Queries) error {
		outboxRepository := outboxRepo.NewOutboxRepository(q)
		outboxService := outboxService.NewOutboxService(outboxRepository)
		log.Println("Fetching task for update...")

		task, err := outboxService.FetchOnePendingForUpdate(ctx)
		if err != nil {
			log.Println("Error fetching task for update:", err)
			return err
		}
		log.Println("Fetched task:", task.ID)
		//orderFullEntity, err := o.orderFullService.GetOrderFull(ctx, task.AggregateID)
		if err != nil {
			log.Println(err)
			return err
		}
		/*
			json, err := json.Marshal(orderFullEntity)
			if err != nil {
				log.Println(err)
				return err
			}

			resp, err := http.Post("localhost:1234", "application/json", bytes.NewBuffer(json))
			defer resp.Body.Close()
			if err != nil {
				log.Println(err)
				return err
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("http status code %d", resp.StatusCode)
			}*/

		if err := outboxService.MarkEventProcessed(ctx, task.ID); err != nil {
			log.Println("Error updating task status in outbox:", err)
			return err
		}
		log.Println(task.ID)
		return nil

		// Задача выполнена успешно в пределах времени
		/*if err := o.outboxService.MarkEventProcessed(ctx, task.ID); err != nil {
			log.Println("Error updating task status in outbox:", err)
			return err

		}
		return nil*/

	})
	log.Println("i am here :", err)
	return err
}

func (o *OrderEno1c) ProcessTasks(ctx context.Context) error {
	//process
	return nil
}

func (o *OrderEno1c) getPendingTasks(ctx context.Context) error {

}

func (o *OrderEno1c) processTask(ctx context.Context) error {
	id := ctx.Value("ID")
	if id == nil {
		return fmt.Errorf("ID not found in context")
	}
	o.transactionalService.ExecTx(ctx, func(q *db.Queries) error {
		//imitate work
		time.Sleep(2 * time.Second)

		outboxRepository := outboxRepo.NewOutboxRepository(q)
		outboxService := outboxService.NewOutboxService(outboxRepository)
		task, err := q.FetchOnePendingForUpdateWithID(ctx, utils.ToUUID(id.(uuid.UUID)))
		if err != nil {
			log.Println("Error fetching task:", err)
		}
		if task.Status == "processed" {
			return nil
		}
		if err := outboxService.MarkEventProcessed(ctx, task.ID.Bytes); err != nil {
			return err
		}

		log.Println(id)
		return nil
	})

	return nil
}
