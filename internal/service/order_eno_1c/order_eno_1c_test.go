package order_eno_1c_test

import (
	"context"
	"orders-center/internal/domain/outbox/entity"
	outboxService "orders-center/internal/domain/outbox/service"
	"orders-center/internal/service/order_eno_1c"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"orders-center/db/sqlc"
)

// MockCron is used to mock the Cron interface
type MockCron struct {
	mock.Mock
}

func (m *MockCron) Stop() {
	m.Called()
}

func (m *MockCron) AddFunc(name string, f func(ctx context.Context) error, interval time.Duration) {
	m.Called(name, f, interval)
}

func (m *MockCron) Start(ctx context.Context) {
	m.Called(ctx)
}

type MockOutboxService struct {
	mock.Mock
}

func (m *MockOutboxService) AddEvent(ctx context.Context, event entity.OutboxEvent) error {
	// Simulate behavior: return nil for success, or an error for failure
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockOutboxService) GetPendingEvents(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]entity.OutboxEvent), args.Error(1)
}

func (m *MockOutboxService) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (entity.OutboxEvent, error) {
	args := m.Called(ctx, eventID, status)
	return args.Get(0).(entity.OutboxEvent), args.Error(1)
}

func (m *MockOutboxService) AddNewEvent(ctx context.Context, event outboxService.AddEventParams) error {
	// Simulate behavior
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockOutboxService) FetchOnePendingForUpdate(ctx context.Context) (entity.OutboxEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).(entity.OutboxEvent), args.Error(1)
}

func (m *MockOutboxService) BatchPendingTasks(ctx context.Context, limit int) ([]entity.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]entity.OutboxEvent), args.Error(1)
}

func (m *MockOutboxService) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOutboxService) FetchOnePendingForUpdateWithID(ctx context.Context, id uuid.UUID) (entity.OutboxEvent, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.OutboxEvent), args.Error(1)
}

func (m *MockOutboxService) MarkEventProcessed(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTransactionService is used to mock the TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) ExecTx(ctx context.Context, fn func(q *db.Queries) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
func TestOrderEno1c(t *testing.T) {
	mockCron := new(MockCron)
	mockOutboxService := new(MockOutboxService)
	mockTransactionService := new(MockTransactionService)

	// Set up the mock responses for the cron jobs
	mockCron.On("AddFunc", "FETCH BATCH", mock.Anything, 1*time.Second).Return()
	mockCron.On("AddFunc", "process task", mock.Anything, 1*time.Second).Return()
	mockCron.On("Start", mock.Anything).Return() // Mock the Start method call

	// Set up the mock responses for the outbox service
	orderID := uuid.Must(uuid.NewV4())
	mockOutboxService.On("BatchPendingTasks", mock.Anything, 10).Return([]entity.OutboxEvent{
		{ID: orderID},
	}, nil)

	mockOutboxService.On("FetchOnePendingForUpdateWithID", mock.Anything, orderID).Return(entity.OutboxEvent{
		ID:          orderID,
		Status:      "pending",
		AggregateID: orderID,
	}, nil)

	mockOutboxService.On("MarkEventProcessed", mock.Anything, orderID).Return(nil)

	// Set up the mock responses for the transaction service
	mockTransactionService.On("ExecTx", mock.Anything, mock.Anything).Return(nil)

	// Create OrderEno1c instance
	orderEno1c := order_eno_1c.NewOrderEno1c(mockCron, mockTransactionService, mockOutboxService)

	// Start cron jobs (this will start the task fetching and processing jobs)
	orderEno1c.Run(context.Background())

	// Assert the mock expectations
	mockOutboxService.AssertExpectations(t)
	mockTransactionService.AssertExpectations(t)
	mockCron.AssertExpectations(t)
}
