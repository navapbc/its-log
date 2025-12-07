-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

-- name: LogIt :one
INSERT INTO itslog (
  version, source, event, value, type
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;