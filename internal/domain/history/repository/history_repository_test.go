package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mockdb "orders-center/db/mock"
	db "orders-center/db/sqlc"
	"testing"
	"time"
)

func TestCreateHistory(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := mockdb.NewMockQuerier(ctrl)
		repo := NewHistoryRepository(mockQuerier)
		oldValue := json.RawMessage(`{"status": "created"}`)
		newValue := json.RawMessage(`{"status": "processing"}`)
		orderID := uuid.Must(uuid.NewV4())
		arg := CreateHistoryParams{
			Type:     "order",
			TypeID:   1,
			OldValue: oldValue,
			Value:    newValue,
			UserID:   "user1",
			OrderID:  orderID,
		}

		expectedHistory := db.History{
			Type:     arg.Type,
			TypeID:   arg.TypeID,
			OldValue: arg.OldValue,
			Value:    arg.Value,
			UserID:   arg.UserID,
			OrderID:  pgtype.UUID{Bytes: orderID, Valid: true},
			Date:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		// Set up mock expectation
		mockQuerier.EXPECT().
			CreateHistory(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, params db.CreateHistoryParams) (db.History, error) {
				assert.Equal(t, arg.Type, params.Type)
				assert.Equal(t, arg.TypeID, params.TypeID)
				assert.Equal(t, arg.Value, params.Value)
				assert.Equal(t, arg.UserID, params.UserID)
				assert.Equal(t, orderID.Bytes(), params.OrderID.Bytes)
				return expectedHistory, nil
			})

		result, err := repo.CreateHistory(context.Background(), arg)

		require.NoError(t, err)
		assert.Equal(t, arg.Type, result.Type)
		assert.Equal(t, arg.TypeID, result.TypeId)
		assert.Equal(t, arg.Value, result.Value)
		assert.Equal(t, arg.UserID, result.UserID)
		assert.Equal(t, orderID, result.OrderID)
	})

	t.Run("error case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockQuerier := mockdb.NewMockQuerier(ctrl)
		repo := NewHistoryRepository(mockQuerier)

		expectedErr := errors.New("database error")
		mockQuerier.EXPECT().
			CreateHistory(gomock.Any(), gomock.Any()).
			Return(db.History{}, expectedErr)

		_, err := repo.CreateHistory(context.Background(), CreateHistoryParams{})
		assert.EqualError(t, err, expectedErr.Error())
	})
}
