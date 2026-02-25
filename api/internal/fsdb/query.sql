-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

--------------------------------------------------------
-- LOGGING
--------------------------------------------------------
-- name: LogEvent :one
INSERT INTO itslog_events (
  timestamp, key_id, cluster_hash, tags_hash, value_hash
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;

-- name: UpdateLookup :exec
INSERT OR IGNORE INTO itslog_lookup (
  timestamp, key_id, hash, name
) VALUES (
  ?, ?, ?, ?
);

--------------------------------------------------------
-- SUMMARY
--------------------------------------------------------
-- name: GetAllSummaries :many
SELECT * FROM itslog_summary;

-- name: InsertSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  date, operation, source_name, event_name, value
  ) VALUES (
  ?, ?, ?, ?, ?
  );

-- name: ReadSummary :many
SELECT 
  date, 
  operation, 
  COALESCE(source_name, '') as source_name, 
  COALESCE(event_name, '') as event_name, 
  value 
FROM itslog_summary
WHERE 
  source_name LIKE COALESCE(?, '%')
  AND
  operation LIKE ?
ORDER BY id
;

--------------------------------------------------------
-- ETL
--------------------------------------------------------

-- name: InsertETL :exec
INSERT OR REPLACE INTO itslog_etl (
  key_id, name, sql
) VALUES (
  ?, ?, ?
);

-- name: GetETL :one
SELECT sql, last_run
FROM itslog_etl
WHERE
  name = ?
LIMIT 1
;

-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = CURRENT_TIMESTAMP 
WHERE 
  name = ?
;


--------------------------------------------------------
-- CLEANUP
--------------------------------------------------------
-- name: VacuumDatabase :exec
VACUUM;

--------------------------------------------------------
-- TEST HELPERS
--------------------------------------------------------



--------------------------------------------------------
-- METADATA
--------------------------------------------------------
-- name: UpdateMeta :exec
-- INSERT OR REPLACE INTO itslog_metadata (
--   key, value
-- ) VALUES (
--   ?, ?
-- );
