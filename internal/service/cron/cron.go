package cron

import (
	"context"
	"log"
	"sync"
	"time"
)

type Cron interface {
	Start(context.Context)
	AddJob(name string, f func(ctx context.Context, taskChan chan<- Task), interval time.Duration)
	AddProcessor(f func(context.Context, Task) error, deadline time.Duration)
	Stop()
}

type Job struct {
	name     string
	job      func(ctx context.Context, taskChan chan<- Task)
	interval time.Duration
}

type Task struct {
	ID   string
	Data interface{}
}

type Scheduler struct {
	jobs        []Job
	processor   func(ctx context.Context, task Task) error
	taskChan    chan Task
	wg          sync.WaitGroup
	workerCount int
	stopChan    chan struct{}
}

func NewScheduler(workerCount int) Cron {
	return &Scheduler{
		taskChan:    make(chan Task, 100),
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

func (s *Scheduler) AddJob(name string, f func(ctx context.Context, taskChan chan<- Task), interval time.Duration) {
	job := Job{
		name:     name,
		job:      f,
		interval: interval,
	}
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) AddProcessor(f func(context.Context, Task) error, deadline time.Duration) {
	s.processor = f
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.scheduleJob(ctx, job)
	}
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}
}

/*func (s *Scheduler) AddFunc(name string, f func(ctx context.Context) error, interval time.Duration, deadline time.Duration) {
	newJob := Job{
		name:     name,
		job:      f,
		interval: interval,
		deadline: deadline,
	}

	s.jobs = append(s.jobs, newJob)
}*/

func (s *Scheduler) scheduleJob(ctx context.Context, job Job) {
	defer s.wg.Done()
	ticker := time.NewTicker(job.interval)
	for {
		select {
		case <-ctx.Done():

			return
		case <-ticker.C:
			log.Println("Starting ", job.name)
			job.job(ctx, s.taskChan)
		}
	}
}

func (s *Scheduler) worker(ctx context.Context, workerID int) {
	defer s.wg.Done()
	for {
		select {
		case task := <-s.taskChan:
			log.Println("worker", workerID, "starts processing task ", task.ID)
			ctxID := context.WithValue(ctx, "id", workerID)
			err := s.processor(ctxID, task)
			if err != nil {
				log.Printf("Worker %d did not compelete processing task %v , err: %v", workerID, task.ID, err.Error())
			}
			log.Println("worker", workerID, "finish task ", task.ID)
		case <-s.stopChan:
			return
		}
	}
}
func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}
