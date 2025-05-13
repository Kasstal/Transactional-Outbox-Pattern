package repository_test

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db "orders-center/db/sqlc"
	"orders-center/internal/domain/history/repository"
	"testing"
	"time"
)

func TestCreateHistory(t *testing.T) {
	// Create a new instance of the mock
	mockQuerier := new(db.MockQuerier)

	// Define the input and expected output
	arg := db.CreateHistoryParams{
		Type:     "order",
		TypeID:   1,
		OldValue: []byte(`{"status":"created"}`),
		Value:    []byte(`{"status":"processing"}`),
		UserID:   "user1",
		OrderID:  pgtype.UUID{Bytes: uuid.Must(uuid.NewV4()), Valid: true},
	}

	expectedHistory := db.History{
		ID:       1,
		Type:     "order",
		TypeID:   1,
		OldValue: []byte(`{"status":"created"}`),
		Value:    []byte(`{"status":"processing"}`),
		Date:     pgtype.Timestamptz(pgtype.Timestamp{Time: time.Now()}),
		UserID:   "user1",
		OrderID:  pgtype.UUID{Bytes: uuid.Must(uuid.NewV4()), Valid: true},
	}

	// Set the expectations on the mock
	mockQuerier.On("CreateHistory", mock.Anything, arg).Return(expectedHistory, nil)

	// Create the repository with the mock
	repo := repository.NewHistoryRepository(mockQuerier)

	// Call the method you want to test

	result, err := repo.CreateHistory(context.Background(), repository.CreateHistoryParams{
		Type:     "order",
		TypeID:   1,
		OldValue: []byte(`{"status":"created"}`),
		Value:    []byte(`{"status":"processing"}`),
		UserID:   "user1",
		OrderID:  uuid.Must(uuid.NewV4()),
	})
	// Assert that no error occurred and the result matches expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedHistory.Type, result.Type)
	assert.Equal(t, expectedHistory.TypeID, result.TypeId)
	assert.Equal(t, expectedHistory.UserID, result.UserID)

	// Verify that the expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestGetHistory(t *testing.T) {
	// Create a new instance of the mock
	mockQuerier := new(db.MockQuerier)

	// Define the input and expected output
	expectedHistory := db.History{
		ID:       1,
		Type:     "order",
		TypeID:   1,
		OldValue: []byte(`{"status":"created"}`),
		Value:    []byte(`{"status":"processing"}`),
		Date:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UserID:   "user1",
		OrderID:  pgtype.UUID{Bytes: uuid.Must(uuid.NewV4()), Valid: true},
	}

	// Set the expectations on the mock
	mockQuerier.On("GetHistory", mock.Anything, int32(1)).Return(expectedHistory, nil)

	// Create the repository with the mock
	repo := repository.NewHistoryRepository(mockQuerier)

	// Call the method you want to test
	result, err := repo.GetHistory(context.Background(), 1)

	// Assert that no error occurred and the result matches expectations
	assert.NoError(t, err)
	assert.Equal(t, expectedHistory.Type, result.Type)
	assert.Equal(t, expectedHistory.TypeID, result.TypeId)
	assert.Equal(t, expectedHistory.UserID, result.UserID)

	// Verify that the expectations were met
	mockQuerier.AssertExpectations(t)
}

func TestDeleteHistory(t *testing.T) {
	// Create a new instance of the mock
	mockQuerier := new(db.MockQuerier)

	// Set the expectations on the mock
	mockQuerier.On("DeleteHistory", mock.Anything, int32(1)).Return(nil)

	// Create the repository with the mock
	repo := repository.NewHistoryRepository(mockQuerier)

	// Call the method you want to test
	err := repo.DeleteHistory(context.Background(), 1)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Verify that the expectations were met
	mockQuerier.AssertExpectations(t)
}
