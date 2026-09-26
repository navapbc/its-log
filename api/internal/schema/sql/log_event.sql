-- name: LogEvent :one
INSERT INTO itslog_events (
  timestamp, key_id, cluster, tags, value
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;
