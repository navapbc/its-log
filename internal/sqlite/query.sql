-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

-- name: LogEvent :one
INSERT INTO events (
  source, event
) VALUES (
  ?, ?
)
RETURNING id;


-- name: UpdateDictionary :exec
INSERT OR IGNORE INTO dictionary (
  event_source, event_name, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?
);