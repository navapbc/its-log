-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

--------------------------------------------------------
-- LOGGING
--------------------------------------------------------
-- name: LogEvent :one
INSERT INTO itslog_events (
  source_hash, event_hash
) VALUES (
  ?, ?
)
RETURNING id;

-- name: LogEventWithValue :one
INSERT INTO itslog_events (
  source_hash, event_hash, value_hash
) VALUES (
  ?, ?, ?
)
RETURNING id;

-- name: LogClusteredEvent :one
INSERT INTO itslog_events (
  cluster_hash, source_hash, event_hash
) VALUES (
  ?, ?, ?
)
RETURNING id;

-- name: LogClusteredEventWithValue :one
INSERT INTO itslog_events (
  timestamp, cluster_hash, source_hash, event_hash, value_hash
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;

-- This is largely for generating fake entries. 
-- However, there may be times where we want to be 
-- more explicit about the timestamp of an entry.
-- name: LogTimestampedEvent :one
INSERT INTO itslog_events (
  timestamp, source_hash, event_hash
) VALUES (
  ?, ?, ?
)
RETURNING id;

-- name: UpdateDictionary :exec
INSERT OR IGNORE INTO itslog_dictionary (
  timestamp, source_name, event_name, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: UpdateLookup :exec
INSERT OR IGNORE INTO itslog_lookup (
  timestamp, hash, name
) VALUES (
  ?, ?, ?
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
-- METADATA
--------------------------------------------------------
-- name: UpdateMeta :exec
INSERT OR REPLACE INTO itslog_metadata (
  key, value
) VALUES (
  ?, ?
);

--------------------------------------------------------
-- ETL
--------------------------------------------------------

-- name: InsertETL :exec
INSERT OR REPLACE INTO itslog_etl (
  name, sql
) VALUES (
  ?, ?
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
-- Used for unit/end-to-end testing.
-- name: TestEventPairExists :one
SELECT EXISTS(
  SELECT 1 
  FROM itslog_events 
  WHERE 
    source_hash = ?
    AND
    event_hash = ?
  );


-- name: TestDictionaryPairExists :one
SELECT EXISTS(
  SELECT 1 
  FROM itslog_dictionary
  WHERE 
    source_hash = ?
    AND
    event_hash = ?
  );
