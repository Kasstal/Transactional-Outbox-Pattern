-- name: CreatePayment :one
INSERT INTO payments (
  order_id, type, sum, info, contract_number, external_id
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments WHERE id = $1 LIMIT 1;

-- name: GetPaymentsByOrder :many
SELECT * FROM payments
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET payed = $2, updated_at = now()
WHERE id = $1
RETURNING *;