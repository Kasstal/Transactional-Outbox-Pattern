-- name: CreateCreditData :one
INSERT INTO credit_data (
  payment_id, bank, type, number_of_months, 
  pay_sum_per_month, broker_id, iin
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetCreditData :one
SELECT * FROM credit_data
WHERE payment_id = $1 LIMIT 1;

-- name: UpdateCreditData :one
UPDATE credit_data
SET 
  pay_sum_per_month = $2,
  number_of_months = $3
WHERE payment_id = $1
RETURNING *;