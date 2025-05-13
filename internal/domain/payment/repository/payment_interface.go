package repository

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"orders-center/internal/domain/payment/entity"
)

type CreatePaymentParams struct {
	OrderID        pgtype.UUID     `json:"order_id"`
	Type           string          `json:"type"`
	Sum            pgtype.Numeric  `json:"sum"`
	Payed          pgtype.Bool     `json:"payed"`
	Info           pgtype.Text     `json:"info"`
	ContractNumber pgtype.Text     `json:"contract_number"`
	ExternalID     pgtype.Text     `json:"external_id"`
	CreditData     json.RawMessage `json:"credit_data"`
	CardData       json.RawMessage `json:"card_data"`
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, arg CreatePaymentParams) (entity.OrderPayment, error)
	GetPayment(ctx context.Context, id uuid.UUID) (entity.OrderPayment, error)
	DeletePayment(ctx context.Context, id uuid.UUID) error
	GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.OrderPayment, error)
}
