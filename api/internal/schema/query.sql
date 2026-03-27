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

--------------------------------------------------------
-- SUMMARY
--------------------------------------------------------
-- name: GetAllSummaries :many
SELECT * FROM itslog_summary;

-- name: InsertSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  date, operation, tags, value
  ) VALUES (
  ?, ?, COALESCE(?, ""), COALESCE(?, "")
  );

-- name: ReadSummary :one
SELECT 
  date, 
  operation, 
  tags, 
  value,
  count
FROM itslog_summary
WHERE 
  tags LIKE ?
  AND
  operation LIKE ?
ORDER BY id
LIMIT 1
;

-- name: ReadSummaries :many
SELECT 
  date, 
  operation, 
  COALESCE(tags, '') as tags, 
  value,
  count
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
