package cron

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockJob struct {
	mock.Mock
}

func (m *MockJob) Execute(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestSchedulerRun(t *testing.T) {
	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Создаем новый планировщик с 2 воркерами
	scheduler := NewScheduler(2)

	// Моковая задача
	mockJob := new(MockJob)
	mockJob.On("Execute", ctx).Return(nil).Once() // Ожидаем, что задача будет выполнена один раз

	// Добавляем задачу в планировщик
	scheduler.AddFunc("test-job", mockJob.Execute, 500*time.Millisecond)

	// Запускаем планировщик
	go scheduler.Start(ctx)

	// Ожидаем выполнения задачи
	time.Sleep(1 * time.Second)

	// Проверяем, что задача была выполнена
	mockJob.AssertExpectations(t)
}

func TestWorkerProcessingJobs(t *testing.T) {
	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Создаем новый планировщик с 1 воркером
	scheduler := NewScheduler(1)

	// Моковая задача
	mockJob := new(MockJob)
	mockJob.On("Execute", ctx).Return(nil).Once() // Ожидаем, что задача будет выполнена один раз

	// Добавляем задачу в планировщик
	scheduler.AddFunc("test-job", mockJob.Execute, 500*time.Millisecond)

	// Запускаем планировщик
	go scheduler.Start(ctx)

	// Ожидаем выполнения задачи
	time.Sleep(1 * time.Second)

	// Проверяем, что задача была выполнена
	mockJob.AssertExpectations(t)
}

func TestSchedulerStop(t *testing.T) {
	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Создаем новый планировщик с 1 воркером
	scheduler := NewScheduler(1)

	// Моковая задача
	mockJob := new(MockJob)
	mockJob.On("Execute", ctx).Return(nil).Once() // Ожидаем, что задача будет выполнена один раз

	// Добавляем задачу в планировщик
	scheduler.AddFunc("test-job", mockJob.Execute, 500*time.Millisecond)

	// Запускаем планировщик
	go scheduler.Start(ctx)

	// Ожидаем выполнения задачи
	time.Sleep(1 * time.Second)

	// Останавливаем планировщик
	scheduler.Stop()

	// Проверяем, что задача была выполнена
	mockJob.AssertExpectations(t)
}
