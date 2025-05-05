-- name: CreateHistory :one
INSERT INTO history (
  type, type_id, old_value, value, user_id, order_id
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetHistory :one
SELECT * FROM history WHERE id = $1 LIMIT 1;

-- name: ListHistoryByOrder :many
SELECT * FROM history 
WHERE order_id = $1 
ORDER BY date DESC
LIMIT $2 OFFSET $3;

-- name: DeleteHistory :exec
DELETE FROM history WHERE id = $1;