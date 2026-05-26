-- name: GetEventByID :one
SELECT * FROM events
WHERE uuid = ? LIMIT 1;

-- name: GetAllEvents :many
SELECT * FROM events
ORDER BY start;

-- name: CreateEvent :exec
INSERT INTO events (
  uuid, summary, location,
  description, start, end,
  created_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateEvent :exec
UPDATE events
SET summary = ?, location = ?,
description = ?, start = ?, end = ?
WHERE uuid = ?;

-- name: DeleteEvent :exec
DELETE FROM events WHERE uuid = ?;
