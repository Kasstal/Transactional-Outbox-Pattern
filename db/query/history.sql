-- name: CreateHistory :one
INSERT INTO history (
    type, type_id, old_value, value, user_id, order_id
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetHistory :one
SELECT * FROM history WHERE id = $1 LIMIT 1;

-- name: UpdateHistory :one
UPDATE history
SET
    type = $2,
    type_id = $3,
    old_value = $4,
    value = $5,
    user_id = $6,
    date = now()
WHERE id = $1
    RETURNING *;

-- name: DeleteHistory :exec
DELETE FROM history WHERE id = $1;

-- name: GetHistoriesByOrderID :many
SELECT * FROM history WHERE order_id = $1;