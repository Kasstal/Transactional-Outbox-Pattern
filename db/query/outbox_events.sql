-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (
    id, aggregate_type, aggregate_id, event_type,
    payload, status, retry_count
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         ) RETURNING *;

-- name: GetOutboxEvent :one
SELECT * FROM outbox_events WHERE id = $1 LIMIT 1;

-- name: GetPendingOutboxEvents :many
SELECT * FROM outbox_events WHERE status = 'pending' LIMIT $1;

-- name: UpdateOutboxEventStatus :one
UPDATE outbox_events
SET
    status = $2,
    retry_count = $3,
    processed_at = CASE WHEN $2 = 'processed' THEN now() ELSE processed_at END,
    payload = $4
WHERE id = $1
    RETURNING *;

-- name: DeleteOutboxEvent :exec
DELETE FROM outbox_events WHERE id = $1;
