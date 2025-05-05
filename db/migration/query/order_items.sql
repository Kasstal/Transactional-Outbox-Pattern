-- name: CreateOrderItem :one
INSERT INTO order_items (
  product_id, external_id, status, base_price, price,
  earned_bonuses, spent_bonuses, gift, owner_id,
  delivery_id, shop_assistant, warehouse, order_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items WHERE id = $1 LIMIT 1;

-- name: ListOrderItems :many
SELECT * FROM order_items 
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderItemStatus :one
UPDATE order_items
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: CalculateOrderTotal :one
SELECT SUM(price) as total FROM order_items 
WHERE order_id = $1;