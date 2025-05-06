package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func createRandomOrderItem(t *testing.T) OrderItem {
	order := createRandomOrder(t)

	arg := CreateOrderItemParams{
		ID:            int32(randomInt(1, 1000)),
		ProductID:     randomProductID(),
		ExternalID:    pgtype.Text{String: randomExternalID(), Valid: true},
		Status:        randomItemStatus(),
		BasePrice:     pgtype.Numeric{Int: big.NewInt(randomInt(1000, 50000)), Exp: -2, Valid: true},
		Price:         pgtype.Numeric{Int: big.NewInt(randomInt(1000, 50000)), Exp: -2, Valid: true},
		EarnedBonuses: pgtype.Numeric{Int: big.NewInt(randomInt(0, 1000)), Exp: -2, Valid: true},
		SpentBonuses:  pgtype.Numeric{Int: big.NewInt(randomInt(0, 1000)), Exp: -2, Valid: true},
		Gift:          pgtype.Bool{Bool: randomBool(), Valid: true},
		OwnerID:       pgtype.Text{String: randomUserID(), Valid: true},
		DeliveryID:    pgtype.Text{String: randomDeliveryID(), Valid: true},
		ShopAssistant: pgtype.Text{String: randomUserID(), Valid: true},
		Warehouse:     pgtype.Text{String: randomWarehouse(), Valid: true},
		OrderID:       pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
	}

	item, err := testQueries.CreateOrderItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	return item
}

func TestCreateOrderItem(t *testing.T) {

	order := createRandomOrder(t)

	arg := CreateOrderItemParams{
		ID:            int32(randomInt(1, 1000)),
		ProductID:     "prod-123",
		ExternalID:    pgtype.Text{String: "ext-123", Valid: true},
		Status:        "reserved",
		BasePrice:     pgtype.Numeric{Int: big.NewInt(5000), Exp: -2, Valid: true},
		Price:         pgtype.Numeric{Int: big.NewInt(4500), Exp: -2, Valid: true},
		EarnedBonuses: pgtype.Numeric{Int: big.NewInt(100), Exp: -2, Valid: true},
		SpentBonuses:  pgtype.Numeric{Int: big.NewInt(500), Exp: -2, Valid: true},
		Gift:          pgtype.Bool{Bool: false, Valid: true},
		OwnerID:       pgtype.Text{String: "user1", Valid: true},
		DeliveryID:    pgtype.Text{String: "delivery-1", Valid: true},
		ShopAssistant: pgtype.Text{String: "assistant-1", Valid: true},
		Warehouse:     pgtype.Text{String: "warehouse-1", Valid: true},
		OrderID:       pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
	}

	item, err := testQueries.CreateOrderItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item)

	require.Equal(t, arg.ID, item.ID)
	require.Equal(t, arg.ProductID, item.ProductID)
	require.Equal(t, arg.ExternalID, item.ExternalID)
	require.Equal(t, arg.Status, item.Status)
	require.Equal(t, arg.BasePrice, item.BasePrice)
	require.Equal(t, arg.Price, item.Price)
	require.Equal(t, arg.EarnedBonuses, item.EarnedBonuses)
	require.Equal(t, arg.SpentBonuses, item.SpentBonuses)
	require.Equal(t, arg.Gift, item.Gift)
	require.Equal(t, arg.OwnerID, item.OwnerID)
	require.Equal(t, arg.DeliveryID, item.DeliveryID)
	require.Equal(t, arg.ShopAssistant, item.ShopAssistant)
	require.Equal(t, arg.Warehouse, item.Warehouse)
	require.Equal(t, arg.OrderID, item.OrderID)

	require.NotZero(t, item.CreatedAt)
	require.NotZero(t, item.UpdatedAt)
}

func TestGetOrderItem(t *testing.T) {

	item1 := createRandomOrderItem(t)
	item2, err := testQueries.GetOrderItem(context.Background(), item1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ID, item2.ID)
	require.Equal(t, item1.ProductID, item2.ProductID)
	require.Equal(t, item1.ExternalID, item2.ExternalID)
	require.Equal(t, item1.Status, item2.Status)
	require.Equal(t, item1.BasePrice, item2.BasePrice)
	require.Equal(t, item1.Price, item2.Price)
	require.Equal(t, item1.EarnedBonuses, item2.EarnedBonuses)
	require.Equal(t, item1.SpentBonuses, item2.SpentBonuses)
	require.Equal(t, item1.Gift, item2.Gift)
	require.Equal(t, item1.OwnerID, item2.OwnerID)
	require.Equal(t, item1.DeliveryID, item2.DeliveryID)
	require.Equal(t, item1.ShopAssistant, item2.ShopAssistant)
	require.Equal(t, item1.Warehouse, item2.Warehouse)
	require.Equal(t, item1.OrderID, item2.OrderID)
	require.WithinDuration(t, item1.CreatedAt.Time, item2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, item1.UpdatedAt.Time, item2.UpdatedAt.Time, time.Second)
}

func TestUpdateOrderItem(t *testing.T) {

	item1 := createRandomOrderItem(t)

	arg := UpdateOrderItemParams{
		ID:            item1.ID,
		ProductID:     "prod-456",
		ExternalID:    pgtype.Text{String: "ext-456", Valid: true},
		Status:        "shipped",
		BasePrice:     pgtype.Numeric{Int: big.NewInt(7000), Exp: -2, Valid: true},
		Price:         pgtype.Numeric{Int: big.NewInt(6300), Exp: -2, Valid: true},
		EarnedBonuses: pgtype.Numeric{Int: big.NewInt(200), Exp: -2, Valid: true},
		SpentBonuses:  pgtype.Numeric{Int: big.NewInt(700), Exp: -2, Valid: true},
		Gift:          pgtype.Bool{Bool: true, Valid: true},
		OwnerID:       pgtype.Text{String: "user2", Valid: true},
		DeliveryID:    pgtype.Text{String: "delivery-2", Valid: true},
		ShopAssistant: pgtype.Text{String: "assistant-2", Valid: true},
		Warehouse:     pgtype.Text{String: "warehouse-2", Valid: true},
	}

	item2, err := testQueries.UpdateOrderItem(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item1.ID, item2.ID)
	require.Equal(t, arg.ProductID, item2.ProductID)
	require.Equal(t, arg.ExternalID, item2.ExternalID)
	require.Equal(t, arg.Status, item2.Status)
	require.Equal(t, arg.BasePrice, item2.BasePrice)
	require.Equal(t, arg.Price, item2.Price)
	require.Equal(t, arg.EarnedBonuses, item2.EarnedBonuses)
	require.Equal(t, arg.SpentBonuses, item2.SpentBonuses)
	require.Equal(t, arg.Gift, item2.Gift)
	require.Equal(t, arg.OwnerID, item2.OwnerID)
	require.Equal(t, arg.DeliveryID, item2.DeliveryID)
	require.Equal(t, arg.ShopAssistant, item2.ShopAssistant)
	require.Equal(t, arg.Warehouse, item2.Warehouse)
	require.Equal(t, item1.OrderID, item2.OrderID)
	require.WithinDuration(t, item1.CreatedAt.Time, item2.CreatedAt.Time, time.Second)
	require.True(t, item2.UpdatedAt.Time.After(item1.UpdatedAt.Time))
}

func TestDeleteOrderItem(t *testing.T) {

	item1 := createRandomOrderItem(t)

	err := testQueries.DeleteOrderItem(context.Background(), item1.ID)
	require.NoError(t, err)

	item2, err := testQueries.GetOrderItem(context.Background(), item1.ID)
	require.Error(t, err)
	require.Empty(t, item2)
}
