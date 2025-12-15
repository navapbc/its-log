-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

-- name: LogIt :one
INSERT INTO itslog (
  source, event
) VALUES (
  ?, ?
)
RETURNING id;