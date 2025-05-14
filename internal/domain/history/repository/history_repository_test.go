package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"orders-center/internal/domain/history/repository"
	"strconv"
	"testing"
	"time"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	mockdb "orders-center/db/mock" // adjust this import path to where your mock is located
	db "orders-center/db/sqlc"
	transactional "orders-center/internal/service/transactional"
)

func setupRepository(t *testing.T) (repository.HistoryRepository, *mockdb.MockQuerier, *pgxpool.Pool) {
	ctrl := gomock.NewController(t)
	mockQuerier := mockdb.NewMockQuerier(ctrl)
	pool := pgxpoolmock.NewMockPgxIface(ctrl)

	repo := repository.NewHistoryRepository(pool)

	return repo, mockQuerier, pool
}

func TestCreateHistory(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())
		now := time.Now()
		dbHistory := db.History{
			ID:       1,
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage("old"),
			Value:    json.RawMessage("new"),
			Date:     pgtype.Timestamptz{Time: now, Valid: true},
			UserID:   "1",
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		// Setup expectations
		expectedArg := db.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage("old"),
			Value:    json.RawMessage("new"),
			UserID:   "1",
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		mockQuerier.EXPECT().
			CreateHistory(gomock.Any(), expectedArg).
			Return(dbHistory, nil)

		arg := repository.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage("old"),
			Value:    json.RawMessage("new"),
			UserID:   "1",
			OrderID:  orderID,
		}

		result, err := repo.CreateHistory(context.Background(), arg)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Type)
		assert.Equal(t, int32(1), result.TypeId)
		assert.Equal(t, "old", result.OldValue)
		assert.Equal(t, "new", result.Value)
		assert.Equal(t, now, result.Date)
		assert.Equal(t, int32(1), result.UserID)
		assert.Equal(t, orderID, result.OrderID)
	})

	t.Run("with transaction", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())
		now := time.Now()
		dbHistory := db.History{
			ID:       1,
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage(`{"status": "processing"}`),
			Value:    json.RawMessage(`{"status": "completed"}`),
			Date:     pgtype.Timestamptz{Time: now, Valid: true},
			UserID:   strconv.Itoa(1),
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		expectedArg := db.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage(`{"status": "processing"}`),
			Value:    json.RawMessage(`{"status": "completed"}`),
			UserID:   strconv.Itoa(1),
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		mockQuerier.EXPECT().
			CreateHistory(gomock.Any(), expectedArg).
			Return(dbHistory, nil)

		arg := repository.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage(`{"status": "processing"}`),
			Value:    json.RawMessage(`{"status": "completed"}`),
			UserID:   strconv.Itoa(1),
			OrderID:  orderID,
		}

		// Create a context with transaction
		mockTx := new(mockdb.MockTx)
		ctx := transactional.NewTxContext(context.Background(), mockTx)

		result, err := repo.CreateHistory(ctx, arg)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Type)
	})

	t.Run("database error", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())
		expectedArg := db.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage(`{"status": "processing"}`),
			Value:    json.RawMessage(`{"status": "completed"}`),
			UserID:   strconv.Itoa(1),
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		mockQuerier.EXPECT().
			CreateHistory(gomock.Any(), expectedArg).
			Return(db.History{}, errors.New("database error"))

		arg := repository.CreateHistoryParams{
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage(`{"status": "processing"}`),
			Value:    json.RawMessage(`{"status": "completed"}`),
			UserID:   strconv.Itoa(1),
			OrderID:  orderID,
		}

		_, err := repo.CreateHistory(context.Background(), arg)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestGetHistory(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())
		now := time.Now()
		dbHistory := db.History{
			ID:       1,
			Type:     "test",
			TypeID:   1,
			OldValue: json.RawMessage("old"),
			Value:    json.RawMessage("new"),
			UserID:   strconv.Itoa(1),
			Date:     pgtype.Timestamptz{Time: now, Valid: true},
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
		}

		mockQuerier.EXPECT().
			GetHistory(gomock.Any(), int32(1)).
			Return(dbHistory, nil)

		result, err := repo.GetHistory(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Type)
		assert.Equal(t, int32(1), result.TypeId)
		assert.Equal(t, now, result.Date)
	})

	t.Run("not found", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		mockQuerier.EXPECT().
			GetHistory(gomock.Any(), int32(1)).
			Return(db.History{}, errors.New("not found"))

		_, err := repo.GetHistory(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})
}

func TestDeleteHistory(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		mockQuerier.EXPECT().
			DeleteHistory(gomock.Any(), int32(1)).
			Return(nil)

		err := repo.DeleteHistory(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		mockQuerier.EXPECT().
			DeleteHistory(gomock.Any(), int32(1)).
			Return(errors.New("database error"))

		err := repo.DeleteHistory(context.Background(), 1)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestGetHistoriesByOrderID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())
		now := time.Now()
		dbHistories := []db.History{
			{
				ID:       1,
				Type:     "test1",
				TypeID:   1,
				OldValue: json.RawMessage("old1"),
				Value:    json.RawMessage("new1"),
				Date:     pgtype.Timestamptz{Time: now, Valid: true},
				UserID:   strconv.Itoa(1),
				OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
			},
			{
				ID:       2,
				Type:     "test2",
				TypeID:   2,
				OldValue: json.RawMessage("old2"),
				Value:    json.RawMessage("new2"),
				Date:     pgtype.Timestamptz{Time: now.Add(time.Hour), Valid: true},
				UserID:   strconv.Itoa(2),
				OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
			},
		}

		mockQuerier.EXPECT().
			GetHistoriesByOrderID(gomock.Any(), pgtype.UUID{Bytes: orderID, Valid: true}).
			Return(dbHistories, nil)

		result, err := repo.GetHistoriesByOrderID(context.Background(), orderID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "test1", result[0].Type)
		assert.Equal(t, "test2", result[1].Type)
	})

	t.Run("no records found", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())

		mockQuerier.EXPECT().
			GetHistoriesByOrderID(gomock.Any(), pgtype.UUID{Bytes: orderID, Valid: true}).
			Return([]db.History{}, nil)

		_, err := repo.GetHistoriesByOrderID(context.Background(), orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no history records")
	})

	t.Run("database error", func(t *testing.T) {
		repo, mockQuerier, _ := setupRepository(t)

		orderID := uuid.Must(uuid.NewV4())

		mockQuerier.EXPECT().
			GetHistoriesByOrderID(gomock.Any(), pgtype.UUID{Bytes: orderID, Valid: true}).
			Return([]db.History{}, errors.New("database error"))

		_, err := repo.GetHistoriesByOrderID(context.Background(), orderID)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}
