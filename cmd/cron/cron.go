package cron

import (
	"context"
	"github.com/gofrs/uuid"
	"log"
	transactional "orders-center/cmd/transactional"
	taskEntity "orders-center/internal/domain/outbox/entity"
	outboxService "orders-center/internal/domain/outbox/service"
	"sync"
	"time"
)

type Scheduler interface {
	Start(ctx context.Context, job func(ctx context.Context) error) // передаем контекст в job
}

type WorkerPoolCron struct {
	display     []uuid.UUID
	outBox      outboxService.OutboxService
	txService   transactional.TransactionService
	interval    time.Duration
	taskChan    chan taskEntity.OutboxEvent
	workerQueue chan struct{}
	workerCount int
	wg          *sync.WaitGroup
}

func NewWorkerPoolScheduler(interval time.Duration, workerCnt int, outboxService outboxService.OutboxService) *WorkerPoolCron {
	return &WorkerPoolCron{
		interval:    interval,
		outBox:      outboxService,
		workerQueue: make(chan struct{}, workerCnt),
		workerCount: workerCnt,
		wg:          &sync.WaitGroup{},
	}
}

/*
	func (w *WorkerPoolCron) Start(ctx context.Context, job func(ctx context.Context) error) {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		// Запуск Cron
		for {
			select {
			case <-ctx.Done():
				return // если контекст отменен, останавливаем cron
			case <-ticker.C:

				for i := len(w.workerQueue); i < w.workerCount; i++ {
					log.Println("workers number: ", len(w.workerQueue))
					w.wg.Add(1)
					w.workerQueue <- struct{}{}
					go func(workerID int) {
						defer w.wg.Done()

						log.Printf("Worker %d started\n", workerID)

						if err := job(ctx); err != nil {
							log.Printf("Error in worker %d: %v\n", workerID, err)
						}

						log.Printf("Worker %d finished\n", workerID)

						<-w.workerQueue
					}(i)

				}

				w.wg.Wait()
			}
		}
	}
*/
func (w *WorkerPoolCron) Start(ctx context.Context, workerJob func(ctx context.Context) error) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	w.taskChan = make(chan taskEntity.OutboxEvent, 10)
	go w.runWorkers(ctx, workerJob)

	// Запуск Cron
	for {
		select {
		case <-ctx.Done():
			w.wg.Wait()
			return // если контекст отменен, останавливаем cron
		case <-ticker.C:
			
		}

		}
	}
}

func worker(ctx context.Context, workerID int, taskChan <-chan taskEntity.OutboxEvent, job func(ctx context.Context) error) {
	for task := range taskChan {
		ctxWithID := context.WithValue(ctx, "ID", task.ID)
		log.Printf("Worker %d starts task with id %v", workerID, task.ID)
		if err := job(ctxWithID); err != nil {
			log.Printf("Error in worker %d: %v\n", workerID, err)
		}
		log.Printf("Worker %d finished\n", workerID)
	}
}

func (w *WorkerPoolCron) runWorkers(ctx context.Context, job func(ctx context.Context) error) {

	for i := 0; i < w.workerCount; i++ {
		go worker(ctx, i, w.taskChan, job)
	}
}
