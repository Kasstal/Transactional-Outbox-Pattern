package db

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomHistory(t *testing.T) History {
	order := createRandomOrder(t)

	oldValue := json.RawMessage(`{"status": "` + randomOrderStatus() + `"}`)
	newValue := json.RawMessage(`{"status": "` + randomOrderStatus() + `"}`)

	arg := CreateHistoryParams{
		Type:     randomHistoryType(),
		TypeID:   int32(randomInt(1, 100)),
		OldValue: oldValue,
		Value:    newValue,
		UserID:   randomUserID(),
		OrderID:  pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
	}

	history, err := testQueries.CreateHistory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, history)

	return history
}
func TestCreateHistory(t *testing.T) {

	order := createRandomOrder(t)

	oldValue := json.RawMessage(`{"status": "created"}`)
	newValue := json.RawMessage(`{"status": "processing"}`)

	arg := CreateHistoryParams{
		Type:     "order",
		TypeID:   1,
		OldValue: oldValue,
		Value:    newValue,
		UserID:   "user1",
		OrderID:  pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
	}

	history, err := testQueries.CreateHistory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, history)

	require.NotZero(t, history.ID)
	require.Equal(t, arg.Type, history.Type)
	require.Equal(t, arg.TypeID, history.TypeID)
	require.Equal(t, arg.OldValue, history.OldValue)
	require.Equal(t, arg.Value, history.Value)
	require.Equal(t, arg.UserID, history.UserID)
	require.Equal(t, arg.OrderID, history.OrderID)
	require.NotZero(t, history.Date)
}

func TestGetHistory(t *testing.T) {

	history1 := createRandomHistory(t)
	history2, err := testQueries.GetHistory(context.Background(), history1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, history2)

	require.Equal(t, history1.ID, history2.ID)
	require.Equal(t, history1.Type, history2.Type)
	require.Equal(t, history1.TypeID, history2.TypeID)
	require.Equal(t, history1.OldValue, history2.OldValue)
	require.Equal(t, history1.Value, history2.Value)
	require.Equal(t, history1.UserID, history2.UserID)
	require.Equal(t, history1.OrderID, history2.OrderID)
	require.WithinDuration(t, history1.Date.Time, history2.Date.Time, time.Second)
}

func TestUpdateHistory(t *testing.T) {

	history1 := createRandomHistory(t)

	oldValue := json.RawMessage(`{"status": "processing"}`)
	newValue := json.RawMessage(`{"status": "completed"}`)

	arg := UpdateHistoryParams{
		ID:       history1.ID,
		Type:     "order",
		TypeID:   2,
		OldValue: oldValue,
		Value:    newValue,
		UserID:   "user2",
	}

	history2, err := testQueries.UpdateHistory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, history2)

	require.Equal(t, history1.ID, history2.ID)
	require.Equal(t, arg.Type, history2.Type)
	require.Equal(t, arg.TypeID, history2.TypeID)
	require.Equal(t, arg.OldValue, history2.OldValue)
	require.Equal(t, arg.Value, history2.Value)
	require.Equal(t, arg.UserID, history2.UserID)
	require.Equal(t, history1.OrderID, history2.OrderID)
	require.True(t, history2.Date.Time.After(history1.Date.Time))
}

func TestDeleteHistory(t *testing.T) {

	history1 := createRandomHistory(t)

	err := testQueries.DeleteHistory(context.Background(), history1.ID)
	require.NoError(t, err)

	history2, err := testQueries.GetHistory(context.Background(), history1.ID)
	require.Error(t, err)
	require.Empty(t, history2)
}
