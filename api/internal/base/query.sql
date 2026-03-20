-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

--------------------------------------------------------
-- LOGGING
--------------------------------------------------------
-- name: LogEvent :one
INSERT INTO itslog_events (
  timestamp, key_id, cluster, tags, value
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;

-- name: UpdateLookup :exec
INSERT OR IGNORE INTO itslog_lookup (
  timestamp, key_id, kind, hash, name
) VALUES (
  ?, ?, ?, ?, ?
);

--------------------------------------------------------
-- SUMMARY
--------------------------------------------------------
-- name: GetAllSummaries :many
SELECT * FROM itslog_summary;

-- name: InsertSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  date, operation, tags, value
  ) VALUES (
  ?, ?, ?, ?
  );

-- name: ReadSummary :many
SELECT 
  date, 
  operation, 
  COALESCE(tags, '') as tags, 
  value 
FROM itslog_summary
WHERE 
  tags LIKE COALESCE(?, '%')
  AND
  operation LIKE ?
ORDER BY id
;

--------------------------------------------------------
-- ETL
--------------------------------------------------------

-- name: InsertETL :exec
INSERT OR REPLACE INTO itslog_etl (
  key_id, name, kind, body
) VALUES (
  ?, ?, ?, ?
);

-- name: GetETL :one
SELECT name, kind, body, last_run
FROM itslog_etl
WHERE
  name = ?
LIMIT 1;

-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = CURRENT_TIMESTAMP 
WHERE 
  name = ?;

-- name: GetDistinctTags :many
SELECT DISTINCT tags
FROM itslog_events;
