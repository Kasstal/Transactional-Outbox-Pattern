package order_eno_1c

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	outbox "orders-center/internal/domain/outbox/entity"
	outboxService "orders-center/internal/domain/outbox/service"
	"orders-center/internal/service/cron"
	orderFull "orders-center/internal/service/order_full/entity"
	orderFullService "orders-center/internal/service/order_full/order_full_service"
	transactional "orders-center/internal/service/transactional"
	"time"
)

type OrderEno1c struct {
	cron                 cron.Cron
	transactionalService transactional.Transactional
	outboxService        outboxService.OutboxService
	orderFullService     orderFullService.OrderFullService
}

func NewOrderEno1c(cron cron.Cron, transactionalService transactional.Transactional, orderFullService orderFullService.OrderFullService, outboxService outboxService.OutboxService) *OrderEno1c {
	return &OrderEno1c{
		outboxService:        outboxService,
		orderFullService:     orderFullService,
		cron:                 cron,
		transactionalService: transactionalService,
	}
}

func (o *OrderEno1c) Run(ctx context.Context) {
	o.cron.AddJob("Fetch batch", o.getPendingTasks, 1*time.Millisecond)
	o.cron.AddProcessor(o.processTask, 6*time.Second)
	o.cron.Start(ctx)
}

func (o *OrderEno1c) getPendingTasks(ctx context.Context, taskChan chan<- cron.Task) {
	err := o.transactionalService.ExecTx(ctx, func(ctx context.Context) error {
		batch, err := o.outboxService.GetPendingEvents(ctx, 2)
		if err != nil {
			log.Print(err)
			return err
		}
		for _, task := range batch {
			if task.ID.IsNil() {
				continue
			}
			log.Println("sending: ", task.ID)
			taskChan <- cron.Task{ID: task.ID.String(), Data: task}
		}
		return nil
	})

	if err != nil {
		log.Printf("failed to fetch pending tasks: %v", err)
	}
}

func (o *OrderEno1c) processTask(ctx context.Context, task cron.Task) error {
	id := ctx.Value("id").(int)
	err := o.transactionalService.ExecTx(ctx, func(ctx context.Context) error {

		outboxTask := task.Data.(outbox.OutboxEvent)
		sqlc, err := o.outboxService.FetchOnePendingForUpdateWithID(ctx, outboxTask.ID)
		if sqlc.Status == "processed" {
			return nil
		}
		log.Printf("Worker %d FETCHED TASK : %v", id, task.ID)
		if err != nil {

			log.Println("Error fetching task:", err, "In Worker: ", id)
			return err
		}
		if err = o.outboxService.MarkEventProcessed(ctx, outboxTask.ID); err != nil {
			return err
		}
		log.Println("marked processed task:", outboxTask.ID, "in Worker: ", id)
		/*if outboxTask.Status == "processed" {
			return fmt.Errorf("task already processed")
		}*/

		orderFull, err := o.orderFullService.GetOrderFull(ctx, outboxTask.AggregateID)
		log.Printf("Order Full : %v in Worker %d", orderFull.Order.ID, id)
		if err != nil {
			log.Println("Could not get OrderFull: ", err)
			//Increment try count
			incrementErr := o.outboxService.IncrementRetryCount(ctx, outboxTask.ID)
			if incrementErr != nil {
				return fmt.Errorf("Could not increment retry count: %v", incrementErr)
			}
			log.Printf("Could not get OrderFull: %v", err, incrementErr)
			return nil
		}

		if err = PostOrderFull(orderFull); err != nil {
			log.Printf("Could not Post OrderFull: %v", err)
		}
		if err != nil {
			incrementErr := o.outboxService.IncrementRetryCount(ctx, outboxTask.ID)
			if incrementErr != nil {
				return fmt.Errorf("Could not increment retry count: %v", incrementErr)
			}
			log.Printf("Could not get OrderFull: %v", err)
			return nil
		}

		log.Println("processed: ", outboxTask.ID, "in Worker: ", id)
		return nil
	})

	return err
}

func (o *OrderEno1c) Stop() error {
	o.cron.Stop()
	return nil
}

func PostOrderFull(orderFull orderFull.OrderFull) error {
	data, err := json.Marshal(orderFull)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "http://localhost:1234/orders", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		log.Printf("failed to send request: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
		log.Printf("server returned non-created status: %v", resp.Status)
	}

	log.Println("Order successfully posted!")
	resp.Body.Close()

	return nil
}
