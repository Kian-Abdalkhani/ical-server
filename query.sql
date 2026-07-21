-- name: GetEventByID :one
SELECT * FROM events
WHERE end > datetime('now') AND uuid = ?
LIMIT 1;

-- name: GetAllEvents :many
SELECT * FROM events
WHERE end > datetime('now')
ORDER BY start;

-- name: CreateEvent :exec
INSERT INTO events (
  uuid, summary, location,
  description, timezone, all_day,
  start, end,
  created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateEvent :exec
UPDATE events
SET summary = ?, location = ?,
description = ?, timezone = ?, all_day = ?, start = ?, end = ?
WHERE uuid = ?;

-- name: DeleteEvent :exec
DELETE FROM events WHERE uuid = ?;
