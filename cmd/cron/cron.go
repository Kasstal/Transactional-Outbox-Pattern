package cron

import (
	"context"
	"log"
	"sync"
	"time"
)

type Cron interface {
	Start(ctx context.Context)
	AddFunc(name string, f func(ctx context.Context) error, interval time.Duration)
}

type Scheduler struct {
	jobs        []Job
	jobChan     chan Job
	wg          sync.WaitGroup
	workerCount int
	stopChan    chan struct{}
}

type Job struct {
	name     string
	job      func(ctx context.Context) error
	interval time.Duration
}

func NewScheduler(workerCount int) Cron {
	return &Scheduler{
		jobChan:     make(chan Job),
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}
func (s *Scheduler) Start(ctx context.Context) {
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}

	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.scheduleJob(ctx, job)
	}

}
func (s *Scheduler) AddFunc(name string, f func(ctx context.Context) error, interval time.Duration) {
	newJob := Job{
		name:     name,
		job:      f,
		interval: interval,
	}

	s.jobs = append(s.jobs, newJob)
}

func (s *Scheduler) scheduleJob(ctx context.Context, job Job) {
	defer s.wg.Done()
	ticker := time.NewTicker(job.interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case s.jobChan <- job:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (s *Scheduler) worker(ctx context.Context, workerID int) {
	defer s.wg.Done()

	for {
		select {
		case job := <-s.jobChan:
			log.Println("worker", workerID, "start job", job.name)
			err := job.job(ctx)
			if err != nil {
				log.Printf("Worker %d did not compelete job %s: %v", workerID, job.name, err.Error())
			}
			log.Println("worker", workerID, "finish job", job.name)
		case <-s.stopChan:
			return
		}
	}
}
func (s *Scheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}
