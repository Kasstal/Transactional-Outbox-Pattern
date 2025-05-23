// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: payments.sql

package db

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (
    order_id, type, sum, payed, info,
    contract_number, external_id, credit_data, card_data
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         ) RETURNING id, order_id, type, sum, payed, info, contract_number, credit_data, external_id, card_data, created_at, updated_at
`

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

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPayment,
		arg.OrderID,
		arg.Type,
		arg.Sum,
		arg.Payed,
		arg.Info,
		arg.ContractNumber,
		arg.ExternalID,
		arg.CreditData,
		arg.CardData,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Type,
		&i.Sum,
		&i.Payed,
		&i.Info,
		&i.ContractNumber,
		&i.CreditData,
		&i.ExternalID,
		&i.CardData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePayment, id)
	return err
}

const getPayment = `-- name: GetPayment :one
SELECT id, order_id, type, sum, payed, info, contract_number, credit_data, external_id, card_data, created_at, updated_at FROM payments WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPayment(ctx context.Context, id pgtype.UUID) (Payment, error) {
	row := q.db.QueryRow(ctx, getPayment, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Type,
		&i.Sum,
		&i.Payed,
		&i.Info,
		&i.ContractNumber,
		&i.CreditData,
		&i.ExternalID,
		&i.CardData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaymentsByOrderID = `-- name: GetPaymentsByOrderID :many
SELECT id, order_id, type, sum, payed, info, contract_number, credit_data, external_id, card_data, created_at, updated_at FROM payments WHERE order_id = $1
`

func (q *Queries) GetPaymentsByOrderID(ctx context.Context, orderID pgtype.UUID) ([]Payment, error) {
	rows, err := q.db.Query(ctx, getPaymentsByOrderID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.Type,
			&i.Sum,
			&i.Payed,
			&i.Info,
			&i.ContractNumber,
			&i.CreditData,
			&i.ExternalID,
			&i.CardData,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePayment = `-- name: UpdatePayment :one
UPDATE payments
SET
    order_id = $2,
    type = $3,
    sum = $4,
    payed = $5,
    info = $6,
    contract_number = $7,
    external_id = $8,
    credit_data = $9,
    card_data = $10,
    updated_at = now()
WHERE id = $1
    RETURNING id, order_id, type, sum, payed, info, contract_number, credit_data, external_id, card_data, created_at, updated_at
`

type UpdatePaymentParams struct {
	ID             pgtype.UUID     `json:"id"`
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

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) (Payment, error) {
	row := q.db.QueryRow(ctx, updatePayment,
		arg.ID,
		arg.OrderID,
		arg.Type,
		arg.Sum,
		arg.Payed,
		arg.Info,
		arg.ContractNumber,
		arg.ExternalID,
		arg.CreditData,
		arg.CardData,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.Type,
		&i.Sum,
		&i.Payed,
		&i.Info,
		&i.ContractNumber,
		&i.CreditData,
		&i.ExternalID,
		&i.CardData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
