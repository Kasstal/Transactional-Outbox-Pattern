package order_eno_1c

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"orders-center/internal/client"
	inboxSvc "orders-center/internal/domain/inbox/service"
	outbox "orders-center/internal/domain/outbox/entity"
	outboxService "orders-center/internal/domain/outbox/service"
	"orders-center/internal/service/cron"
	orderFullService "orders-center/internal/service/order_full/order_full_service"
	transactional "orders-center/internal/service/transactional"
	"orders-center/internal/utils"
	"time"
)

type OrderEno1c struct {
	cron                 cron.Cron
	transactionalService transactional.Transactional
	outboxService        outboxService.OutboxService
	inboxService         inboxSvc.InboxService
	orderFullService     orderFullService.OrderFullService
	client               *client.Client
	retryMax             int32
	pollerInterval       time.Duration
	workerTimeout        time.Duration
}

func NewOrderEno1c(
	cron cron.Cron,
	transactionalService transactional.Transactional,
	orderFullService orderFullService.OrderFullService,
	outboxService outboxService.OutboxService,
	inboxService inboxSvc.InboxService,
	client *client.Client,
	config utils.Config,
) *OrderEno1c {
	return &OrderEno1c{
		client:               client,
		retryMax:             config.MaxRetries,
		pollerInterval:       config.JobInterval,
		workerTimeout:        config.WorkerTimeout,
		inboxService:         inboxService,
		outboxService:        outboxService,
		orderFullService:     orderFullService,
		cron:                 cron,
		transactionalService: transactionalService,
	}
}

func (o *OrderEno1c) Run(ctx context.Context) {
	o.cron.AddJob("Fetch batch", o.getPendingTasks, o.pollerInterval)
	o.cron.AddProcessor(o.processTask, o.workerTimeout)
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
		//RETRIEVE TASK
		outboxTask := task.Data.(outbox.OutboxEvent)
		//GET LOCK NOWAIT
		sqlc, err := o.outboxService.FetchOnePendingForUpdateWithID(ctx, outboxTask.ID)

		log.Printf("Worker %d FETCHED TASK : %v", id, task.ID)
		if err != nil {

			log.Println("Error fetching task:", err, "In Worker: ", id)
			return err
		}
		//CHECK IF THIS EVENT WAS ALREADY PROCESSED AND SENT
		if o.checkProcessed(ctx, sqlc.ID) {
			return fmt.Errorf("task id %d is already processed In Worker: %d", outboxTask.ID, id)
		}

		//MARK EVENT PROCESSED SO THAT NEXT BATCH WILL NOT POLL IT
		if err = o.outboxService.MarkEventProcessed(ctx, outboxTask.ID); err != nil {
			return err
		}
		//RETRIEVE ORDER FULL
		orderFull, err := o.orderFullService.GetOrderFull(ctx, outboxTask.AggregateID)
		log.Printf("Order Full : %v In Worker %d", orderFull.Order.ID, id)
		if err != nil {
			log.Println("Could not get OrderFull: ", err)
			//Increment try count
			incrementErr := o.incrementOrFail(ctx, outboxTask.ID, err.Error())
			if incrementErr != nil {
				return fmt.Errorf("could not get OrderFull : %v Increment err: %v In Worker: %d", err, incrementErr, id)
			}
			return nil
		}

		//SENDING TO MOCK1C
		resp, err := o.client.SendRequest("orders", "POST", orderFull)
		if resp == nil {
			incrementErr := o.incrementOrFail(ctx, outboxTask.ID, err.Error())
			if incrementErr != nil {
				return fmt.Errorf("could POST OrderFull : %v Increment err: %v In Worker: %d", err, incrementErr, id)
			}
			return nil
		}
		if resp.StatusCode != http.StatusCreated || err != nil {
			log.Printf("Could not Post OrderFull: %v in Worker:%d", err, id)
			incrementErr := o.incrementOrFail(ctx, outboxTask.ID, err.Error())
			if incrementErr != nil {
				return fmt.Errorf("could POST OrderFull : %v Increment err: %v In Worker: %d", err, incrementErr, id)
			}

			return nil
		}
		//CREATE NEW TASK PROCESSED RECORD IN INBOX
		_, err = o.inboxService.Create(ctx, sqlc.ID)
		if err != nil {
			return err
		}
		log.Println("added processed task id into INBOX: ", sqlc.ID)
		log.Println("processed: ", outboxTask.ID, "in Worker: ", id)
		return nil
	})

	return err
}

func (o *OrderEno1c) incrementOrFail(ctx context.Context, id uuid.UUID, errMsg string) error {
	retryCount, incrementErr := o.outboxService.IncrementRetryCount(ctx, id, errMsg)
	if incrementErr != nil {
		return fmt.Errorf("could not increment retry count: %v", incrementErr)
	}

	//failed
	if retryCount >= o.retryMax {
		err := o.outboxService.MarkFailed(ctx, id, errMsg)
		if err != nil {
			return fmt.Errorf("could not update event status to failed: %v", err)
		}
	}
	return nil
}

func (o *OrderEno1c) Stop() error {
	o.cron.Stop()
	return nil
}

func (o *OrderEno1c) CreateOutboxTask(ctx context.Context, orderID uuid.UUID) error {

	eventData := map[string]interface{}{
		"order_id": orderID,
	}
	payload, err := json.Marshal(eventData)
	if err != nil {
		log.Println(err)
		return err
	}
	// Add event в Outbox
	if err = o.outboxService.AddNewEvent(ctx, outboxService.AddEventParams{
		AggregateType: "OrderFull",
		AggregateID:   orderID,
		EventType:     "OrderCreated",
		Payload:       payload,
	}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (o *OrderEno1c) checkProcessed(ctx context.Context, id uuid.UUID) bool {
	getEventIfExists, err := o.inboxService.GetInboxEvent(ctx, id)
	if err != nil {
		log.Printf("could not get inbox event: %v", err)
	}
	if getEventIfExists.EventID.IsNil() {
		return false
	}
	return true
}

/*func PostOrderFull(orderFull orderFullEntity.OrderFull) error {
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
		log.Printf("server returned non-created status: %v", resp.Status)
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)

	}

	log.Println("Order successfully posted!")
	resp.Body.Close()

	return nil
}*/
