-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

--------------------------------------------------------
-- LOGGING
--------------------------------------------------------
-- name: LogEvent :one
INSERT INTO itslog_events (
  key_id, source_hash, event_hash
) VALUES (
  ?, ?, ?
)
RETURNING id;

-- name: LogEventWithValue :one
INSERT INTO itslog_events (
  key_id, source_hash, event_hash, value_hash
) VALUES (
  ?, ?, ?, ?
)
RETURNING id;

-- name: LogClusteredEvent :one
INSERT INTO itslog_events (
  key_id, cluster_hash, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?
)
RETURNING id;

-- name: LogClusteredEventWithValue :one
INSERT INTO itslog_events (
  key_id, timestamp, cluster_hash, source_hash, event_hash, value_hash
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- This is largely for generating fake entries. 
-- However, there may be times where we want to be 
-- more explicit about the timestamp of an entry.
-- name: LogTimestampedEvent :one
INSERT INTO itslog_events (
  key_id, timestamp, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?
)
RETURNING id;

-- name: UpdateDictionary :exec
INSERT OR IGNORE INTO itslog_dictionary (
  key_id, source_name, event_name, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: UpdateLookup :exec
INSERT OR IGNORE INTO itslog_lookup (
  key_id, hash, name
) VALUES (
  ?, ?, ?
);

--------------------------------------------------------
-- METADATA
--------------------------------------------------------
-- name: UpdateMeta :exec
INSERT OR REPLACE INTO itslog_metadata (
  key_id, key, value
) VALUES (
  ?, ?, ?
);

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
  key_id = ?
  AND
  name = ?
LIMIT 1
;

-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = CURRENT_TIMESTAMP 
WHERE 
  key_id = ?
  AND
  name = ?
;


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
