-- name: GetProcessedEvent :one
SELECT * FROM inbox_events WHERE event_id = $1 LIMIT 1;


-- name: CreateInboxEvent :one
INSERT INTO inbox_events (
    event_id
) VALUES (
             $1
         ) RETURNING *;