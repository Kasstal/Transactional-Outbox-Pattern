package order_eno_1c

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"orders-center/internal/service/cron"
	orderFull "orders-center/internal/service/order_full/entity"
	orderFullService "orders-center/internal/service/order_full/order_full_service"
	transactional "orders-center/internal/service/transactional"

	outbox "orders-center/internal/domain/outbox/entity"
	outboxService "orders-center/internal/domain/outbox/service"
	"sync"
	"time"
)

type OrderEno1c struct {
	mu                   sync.Mutex
	cron                 cron.Cron
	tasks                []outbox.OutboxEvent
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
	o.cron.AddFunc("FETCH BATCH", o.getPendingTasks, 1*time.Second, 3*time.Second)
	o.cron.AddFunc("process task", o.processTask, 200*time.Millisecond, 1*time.Second)
	o.cron.Start(ctx)
}

func (o *OrderEno1c) getPendingTasks(ctx context.Context) error {
	batch, err := o.outboxService.BatchPendingTasks(ctx, 15)
	if err != nil {
		return err
	}
	for _, task := range batch {

		if task.ID.String() == "00000000-0000-0000-0000-000000000000" {
			continue
		}
		log.Println("sending: ", task.ID)
		o.mu.Lock()
		o.tasks = append(o.tasks, task)
		o.mu.Unlock()
	}

	return nil
}

func (o *OrderEno1c) processTask(ctx context.Context) error {
	if len(o.tasks) == 0 {
		return fmt.Errorf("no pending tasks")
	}
	o.mu.Lock()
	id := o.tasks[0].ID
	log.Println("processing task: ", id)
	o.mu.Unlock()
	o.transactionalService.ExecTx(ctx, func(ctx context.Context) error {
		//imitate work

		task, err := o.outboxService.FetchOnePendingForUpdateWithID(ctx, id)
		log.Printf("FETCHED TASK : %v", task.ID)
		if err != nil {
			log.Println("Error fetching task:", err)
			return err
		}
		if task.Status == "processed" {
			return nil
		}

		orderFull, err := o.orderFullService.GetOrderFull(ctx, task.AggregateID)
		log.Printf("Order Full : %v", orderFull.Order.ID)
		if err != nil {
			log.Println("Could not get OrderFull: ", err)
		}

		if err = PostOrderFull(orderFull); err != nil {
			log.Printf("Could not Post OrderFull: %v", err)
		}
		if err != nil {
			incrementErr := o.outboxService.IncrementRetryCount(ctx, task.ID)
			if incrementErr != nil {
				return fmt.Errorf("Could not increment retry count: %v", incrementErr)
			}
			log.Printf("Could not get OrderFull: %v Could not increment: %v", err, incrementErr)
			return nil
		}
		if err = o.outboxService.MarkEventProcessed(ctx, task.ID); err != nil {
			return err
		}

		log.Println("processed: ", id)
		return nil
	})
	o.mu.Lock()
	o.tasks = o.tasks[1:]
	o.mu.Unlock()
	return nil
}

func (o *OrderEno1c) InitReset() error {
	ctx := context.Background()
	inProgressTasks, err := o.outboxService.GetAllInProgressEvents(ctx)
	for _, task := range inProgressTasks {
		_, err = o.outboxService.UpdateEventStatus(context.Background(), task.ID, "pending")
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderEno1c) ShutdownReset() error {
	var err error

	for _, task := range o.tasks {
		_, err = o.outboxService.UpdateEventStatus(context.Background(), task.ID, "pending")
		if err != nil {
			return err
		}
	}
	return nil
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
