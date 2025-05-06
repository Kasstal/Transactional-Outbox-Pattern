package db

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func createRandomPayment(t *testing.T) Payment {
	order := createRandomOrder(t)

	creditData := json.RawMessage(`{"bank": "` + randomBank() + `", "term": "` + strconv.Itoa(int(randomInt(1, 36))) + `"}`)
	cardData := json.RawMessage(`{"last4": "` + randomLast4() + `", "system": "` + randomCardSystem() + `"}`)
	if !json.Valid(creditData) {
		t.Fatalf("Invalid creditData JSON: %s", creditData)
	}
	if !json.Valid(cardData) {
		t.Fatalf("Invalid cardData JSON: %s", cardData)
	}
	arg := CreatePaymentParams{
		ID:             pgtype.UUID{Bytes: randomUUID(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
		Type:           randomPaymentType(),
		Sum:            pgtype.Numeric{Int: big.NewInt(randomInt(1000, 100000)), Exp: -2, Valid: true},
		Payed:          pgtype.Bool{Bool: randomBool(), Valid: true},
		Info:           pgtype.Text{String: randomPaymentInfo(), Valid: true},
		ContractNumber: pgtype.Text{String: randomContractNumber(), Valid: true},
		ExternalID:     pgtype.Text{String: randomExternalID(), Valid: true},
		CreditData:     creditData,
		CardData:       cardData,
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	return payment
}
func TestCreatePayment(t *testing.T) {

	order := createRandomOrder(t)

	creditData := json.RawMessage(`{"bank": "sberbank", "term": 12}`)
	cardData := json.RawMessage(`{"last4": "1234", "system": "visa"}`)

	arg := CreatePaymentParams{
		ID:             pgtype.UUID{Bytes: randomUUID(), Valid: true},
		OrderID:        pgtype.UUID{Bytes: order.ID.Bytes, Valid: true},
		Type:           "credit",
		Sum:            pgtype.Numeric{Int: big.NewInt(10000), Exp: -2, Valid: true},
		Payed:          pgtype.Bool{Bool: true, Valid: true},
		Info:           pgtype.Text{String: "test payment", Valid: true},
		ContractNumber: pgtype.Text{String: "CNT-123", Valid: true},
		ExternalID:     pgtype.Text{String: "EXT-123", Valid: true},
		CreditData:     creditData,
		CardData:       cardData,
	}

	payment, err := testQueries.CreatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.ID, payment.ID)
	require.Equal(t, arg.OrderID, payment.OrderID)
	require.Equal(t, arg.Type, payment.Type)
	require.Equal(t, arg.Sum, payment.Sum)
	require.Equal(t, arg.Payed, payment.Payed)
	require.Equal(t, arg.Info, payment.Info)
	require.Equal(t, arg.ContractNumber, payment.ContractNumber)
	require.Equal(t, arg.ExternalID, payment.ExternalID)
	require.Equal(t, arg.CreditData, payment.CreditData)
	require.Equal(t, arg.CardData, payment.CardData)

	require.NotZero(t, payment.CreatedAt)
	require.NotZero(t, payment.UpdatedAt)
}

func TestGetPayment(t *testing.T) {

	payment1 := createRandomPayment(t)
	payment2, err := testQueries.GetPayment(context.Background(), payment1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, payment1.OrderID, payment2.OrderID)
	require.Equal(t, payment1.Type, payment2.Type)
	require.Equal(t, payment1.Sum, payment2.Sum)
	require.Equal(t, payment1.Payed, payment2.Payed)
	require.Equal(t, payment1.Info, payment2.Info)
	require.Equal(t, payment1.ContractNumber, payment2.ContractNumber)
	require.Equal(t, payment1.ExternalID, payment2.ExternalID)
	require.Equal(t, payment1.CreditData, payment2.CreditData)
	require.Equal(t, payment1.CardData, payment2.CardData)
	require.WithinDuration(t, payment1.CreatedAt.Time, payment2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, payment1.UpdatedAt.Time, payment2.UpdatedAt.Time, time.Second)
}

func TestUpdatePayment(t *testing.T) {

	payment1 := createRandomPayment(t)

	creditData := json.RawMessage(`{"bank": "tinkoff", "term": 6}`)
	cardData := json.RawMessage(`{"last4": "5678", "system": "mastercard"}`)

	arg := UpdatePaymentParams{
		ID:             payment1.ID,
		OrderID:        payment1.OrderID,
		Type:           "card",
		Sum:            pgtype.Numeric{Int: big.NewInt(15000), Exp: -2, Valid: true},
		Payed:          pgtype.Bool{Bool: false, Valid: true},
		Info:           pgtype.Text{String: "updated payment", Valid: true},
		ContractNumber: pgtype.Text{String: "CNT-456", Valid: true},
		ExternalID:     pgtype.Text{String: "EXT-456", Valid: true},
		CreditData:     creditData,
		CardData:       cardData,
	}

	payment2, err := testQueries.UpdatePayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, arg.OrderID, payment2.OrderID)
	require.Equal(t, arg.Type, payment2.Type)
	require.Equal(t, arg.Sum, payment2.Sum)
	require.Equal(t, arg.Payed, payment2.Payed)
	require.Equal(t, arg.Info, payment2.Info)
	require.Equal(t, arg.ContractNumber, payment2.ContractNumber)
	require.Equal(t, arg.ExternalID, payment2.ExternalID)
	require.Equal(t, arg.CreditData, payment2.CreditData)
	require.Equal(t, arg.CardData, payment2.CardData)
	require.WithinDuration(t, payment1.CreatedAt.Time, payment2.CreatedAt.Time, time.Second)
	require.True(t, payment2.UpdatedAt.Time.After(payment1.UpdatedAt.Time))
}

func TestDeletePayment(t *testing.T) {

	payment1 := createRandomPayment(t)

	err := testQueries.DeletePayment(context.Background(), payment1.ID)
	require.NoError(t, err)

	payment2, err := testQueries.GetPayment(context.Background(), payment1.ID)
	require.Error(t, err)
	require.Empty(t, payment2)
}
