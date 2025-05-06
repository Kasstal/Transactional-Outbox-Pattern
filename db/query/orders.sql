-- name: CreateOrder :one
INSERT INTO orders (
    id, type, status, city, subdivision, price,
    platform, general_id, order_number, executor
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         ) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 LIMIT 1;

-- name: UpdateOrder :one
UPDATE orders
SET
    type = $2,
    status = $3,
    city = $4,
    subdivision = $5,
    price = $6,
    platform = $7,
    general_id = $8,
    order_number = $9,
    executor = $10,
    updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = $1;