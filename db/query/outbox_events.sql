-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (
    aggregate_type, aggregate_id, event_type,
    payload, status, retry_count
) VALUES (
             $1, $2, $3, $4, $5, $6
         ) RETURNING *;

-- name: GetOutboxEvent :one
SELECT * FROM outbox_events WHERE id = $1 LIMIT 1;

-- name: GetPendingOutboxEvents :many
SELECT * FROM outbox_events WHERE status = 'pending'
                            FOR UPDATE SKIP LOCKED
                             LIMIT $1;
-- name: GetAllInProgressOutboxEvents :many
SELECT * FROM outbox_events WHERE status = 'in_progress';

-- name: UpdateOutboxEventStatus :one
UPDATE outbox_events
SET
    status = @status::varchar(20),
    processed_at = CASE WHEN @status = 'processed' THEN now() ELSE processed_at END
WHERE id = sqlc.arg(id)
    RETURNING *;

-- name: DeleteOutboxEvent :exec
DELETE FROM outbox_events WHERE id = $1;


-- name: FetchOnePendingForUpdate :one
SELECT *
FROM outbox_events
WHERE status = 'pending' FOR UPDATE NOWAIT LIMIT 1;

-- name: FetchOnePendingForUpdateWithID :one
SELECT *
FROM outbox_events
WHERE id = $1 FOR UPDATE NOWAIT LIMIT 1;

-- name: BatchPendingTasks :many
WITH batch AS (
    SELECT id
    FROM outbox_events
    WHERE status = 'pending'
    LIMIT $1
    )
UPDATE outbox_events
SET status = 'in_progress'
WHERE id IN (SELECT id FROM batch)
RETURNING *;


-- name: MarkEventError :exec
UPDATE outbox_events
SET
    status = 'failed',
    error_message = $2,
    processed_at = NOW()
WHERE id = $1;


-- name: IncrementRetryCount :one
UPDATE outbox_events
SET retry_count = retry_count + 1,
    status = 'pending',
    error_message = $2
WHERE id = $1 RETURNING retry_count;