package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"math/big"
	"orders-center/internal/utils"
	"testing"
	"time"
)

func createRandomOrder(t *testing.T) Order {
	arg := CreateOrderParams{
		ID:          pgtype.UUID{Bytes: randomUUID(), Valid: true},
		Type:        randomOrderType(),
		Status:      randomOrderStatus(),
		City:        randomCity(),
		Subdivision: pgtype.Text{String: randomSubdivision(), Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(randomInt(1000, 100000)), Exp: -2, Valid: true},
		Platform:    randomPlatform(),
		GeneralID:   pgtype.UUID{Bytes: randomUUID(), Valid: true},
		OrderNumber: utils.ToText(randomOrderNumber()),
		Executor:    pgtype.Text{String: randomUserID(), Valid: true},
	}

	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	return order
}

func TestCreateOrder(t *testing.T) {

	arg := CreateOrderParams{
		ID:          pgtype.UUID{Bytes: randomUUID(), Valid: true},
		Type:        "online",
		Status:      "created",
		City:        "Moscow",
		Subdivision: pgtype.Text{String: "north", Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(10000), Exp: -2, Valid: true},
		Platform:    "web",
		GeneralID:   pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true},
		OrderNumber: utils.ToText("ORD-123"),
		Executor:    pgtype.Text{String: "user1", Valid: true},
	}

	order, err := testQueries.CreateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order)

	require.Equal(t, arg.ID, order.ID)
	require.Equal(t, arg.Type, order.Type)
	require.Equal(t, arg.Status, order.Status)
	require.Equal(t, arg.City, order.City)
	require.Equal(t, arg.Subdivision, order.Subdivision)
	require.Equal(t, arg.Price, order.Price)
	require.Equal(t, arg.Platform, order.Platform)
	require.Equal(t, arg.GeneralID, order.GeneralID)
	require.Equal(t, arg.OrderNumber, order.OrderNumber)
	require.Equal(t, arg.Executor, order.Executor)

	require.NotZero(t, order.CreatedAt)
	require.NotZero(t, order.UpdatedAt)
}

func TestGetOrder(t *testing.T) {

	order1 := createRandomOrder(t)
	order2, err := testQueries.GetOrder(context.Background(), order1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, order1.Type, order2.Type)
	require.Equal(t, order1.Status, order2.Status)
	require.Equal(t, order1.City, order2.City)
	require.Equal(t, order1.Subdivision, order2.Subdivision)
	require.Equal(t, order1.Price, order2.Price)
	require.Equal(t, order1.Platform, order2.Platform)
	require.Equal(t, order1.GeneralID, order2.GeneralID)
	require.Equal(t, order1.OrderNumber, order2.OrderNumber)
	require.Equal(t, order1.Executor, order2.Executor)
	require.WithinDuration(t, order1.CreatedAt.Time, order2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, order1.UpdatedAt.Time, order2.UpdatedAt.Time, time.Second)
}

func TestUpdateOrder(t *testing.T) {

	order1 := createRandomOrder(t)

	arg := UpdateOrderParams{
		ID:          order1.ID,
		Type:        "offline",
		Status:      "processing",
		City:        "Saint Petersburg",
		Subdivision: pgtype.Text{String: "center", Valid: true},
		Price:       pgtype.Numeric{Int: big.NewInt(15000), Exp: -2, Valid: true},
		Platform:    "mobile",
		GeneralID:   pgtype.UUID{Bytes: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, Valid: true},
		OrderNumber: utils.ToText("ORD-456"),
		Executor:    pgtype.Text{String: "user2", Valid: true},
	}

	order2, err := testQueries.UpdateOrder(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, order2)

	require.Equal(t, order1.ID, order2.ID)
	require.Equal(t, arg.Type, order2.Type)
	require.Equal(t, arg.Status, order2.Status)
	require.Equal(t, arg.City, order2.City)
	require.Equal(t, arg.Subdivision, order2.Subdivision)
	require.Equal(t, arg.Price, order2.Price)
	require.Equal(t, arg.Platform, order2.Platform)
	require.Equal(t, arg.GeneralID, order2.GeneralID)
	require.Equal(t, arg.OrderNumber, order2.OrderNumber)
	require.Equal(t, arg.Executor, order2.Executor)
	require.WithinDuration(t, order1.CreatedAt.Time, order2.CreatedAt.Time, time.Second)

}

func TestDeleteOrder(t *testing.T) {

	order1 := createRandomOrder(t)
	err := testQueries.DeleteOrder(context.Background(), order1.ID)
	require.NoError(t, err)

	order2, err := testQueries.GetOrder(context.Background(), order1.ID)
	require.Error(t, err)
	require.Empty(t, order2)
}
