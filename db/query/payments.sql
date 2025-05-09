-- name: CreatePayment :one
INSERT INTO payments (
    id, order_id, type, sum, payed, info,
    contract_number, external_id, credit_data, card_data
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         ) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

-- name: UpdatePayment :one
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
    RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments WHERE id = $1;

-- name: GetPaymentsByOrderID :many
SELECT * FROM payments WHERE order_id = $1;