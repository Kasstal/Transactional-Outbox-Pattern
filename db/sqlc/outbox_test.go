package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"math/rand"
	"orders-center/internal/utils"
	"testing"
)

func createRandomOutboxEvent(t *testing.T) OutboxEvent {

	event := CreateOutboxEventParams{
		AggregateType: randomString("order"),
		AggregateID:   utils.ToUUID(randomUUID()),
		EventType:     randomString("created"),
		Payload:       []byte(`{"order_number": "` + randomString("ORD-") + `"}`),
		Status:        "pending",
		RetryCount:    pgtype.Int4{Int32: int32(rand.Intn(10)), Valid: true},
	}

	// Insert the event into the database
	insertedEvent, err := testQueries.CreateOutboxEvent(context.Background(), event)
	require.NoError(t, err)
	require.NotEmpty(t, insertedEvent)

	return insertedEvent
}

func TestCreateOutboxEvent(t *testing.T) {

	// Create and insert a random outbox event
	insertedEvent := createRandomOutboxEvent(t)

	// Verify the event is inserted correctly
	require.NotZero(t, insertedEvent.ID)
	require.Equal(t, "pending", insertedEvent.Status)
	require.NotNil(t, insertedEvent.CreatedAt)
	require.NotNil(t, insertedEvent.ProcessedAt)
}
