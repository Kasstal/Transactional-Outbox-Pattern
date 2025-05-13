-- name: CreateOrderItem :one
INSERT INTO order_items (
    product_id, external_id, status, base_price,
    price, earned_bonuses, spent_bonuses, gift,
    owner_id, delivery_id, shop_assistant, warehouse, order_id
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
         ) RETURNING *;

-- name: GetOrderItem :one
SELECT * FROM order_items WHERE id = $1 LIMIT 1;

-- name: UpdateOrderItem :one
UPDATE order_items
SET
    product_id = $2,
    external_id = $3,
    status = $4,
    base_price = $5,
    price = $6,
    earned_bonuses = $7,
    spent_bonuses = $8,
    gift = $9,
    owner_id = $10,
    delivery_id = $11,
    shop_assistant = $12,
    warehouse = $13,
    updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1;