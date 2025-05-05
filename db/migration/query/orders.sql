-- name: CreateOrder :one
INSERT INTO orders (
  type, status, city, subdivision, price, platform,
  general_id, order_number, executor
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: GetOrderForUpdate :one
SELECT * FROM orders WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListOrders :many
SELECT * FROM orders 
WHERE executor = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2, updated_at = now() 
WHERE id = $1 
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;