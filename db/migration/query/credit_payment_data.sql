-- name: CreateCardPaymentData :one
INSERT INTO card_payment_data (
  payment_id, provider, transaction_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetCardPaymentData :one
SELECT * FROM card_payment_data
WHERE payment_id = $1 LIMIT 1;

-- name: UpdateCardPaymentData :one
UPDATE card_payment_data
SET 
  transaction_id = $2
WHERE payment_id = $1
RETURNING *;