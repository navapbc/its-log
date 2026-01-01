-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-sqlite.html

-- name: LogEvent :one
INSERT INTO itslog_events (
  source_hash, event_hash
) VALUES (
  ?, ?
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
  source_name, event_name, source_hash, event_hash
) VALUES (
  ?, ?, ?, ?
);

--------------------------------------------------------
-- METADATA
--------------------------------------------------------

-- name: InsertMetadata :exec
INSERT OR REPLACE INTO itslog_metadata (key, value) 
  VALUES (?, ?);

-- name: GetMetadata :one
SELECT key, value FROM itslog_metadata
  WHERE key = ? LIMIT 1;

--------------------------------------------------------
-- SUMMARIZING DATA
--------------------------------------------------------
-- name: InsertSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  operation, source_name, event_name, value 
  ) VALUES (
  ?, ?, ?, ?
  );

------------------------
-- By the hour
------------------------
-- name: EventCountsByTheHour :many
SELECT 
  strftime('%H', timestamp) AS hour,
  source_hash,
  event_hash,
  COUNT(*) AS event_count
FROM itslog_events
GROUP BY hour, source_hash, event_hash
ORDER BY hour, source_hash, event_hash;

-- name: SourceCountsByTheHour :many
SELECT 
  strftime('%H', timestamp) AS hour,
  source_hash,
  event_hash,
  COUNT(*) AS source_count
FROM itslog_events
GROUP BY hour, source_hash
ORDER BY hour, source_hash;

------------------------
-- By the day
------------------------
-- name: EventCountsForTheDay :many
SELECT 
  source_hash,
  event_hash,
  COUNT(*) AS event_count
FROM itslog_events
GROUP BY source_hash, event_hash
ORDER BY source_hash, event_hash;

-- name: SourceCountsForTheDay :many
SELECT 
  source_hash,
  COUNT(*) AS source_count
FROM itslog_events
GROUP BY source_hash
ORDER BY source_hash;

------------------------
-- Summary helpers
------------------------
-- name: GetSourceName :one
SELECT
  source_name
  FROM itslog_dictionary
  WHERE
    source_hash = ?
  LIMIT 1;

-- name: GetEventName :one
SELECT
  event_name
  FROM itslog_dictionary
  WHERE
    source_hash = ?
    AND
    event_hash = ?
  LIMIT 1;

-- name: GetSourceNames :many
SELECT
  source_name
  FROM 
  itslog_dictionary
;

-- name: GetEventNamesForSource :many
SELECT
  event_name
  FROM
  itslog_dictionary
  WHERE 
  source_name = ?
;


-- name: DeleteSummaryData :exec
DELETE FROM itslog_summary;

-- NAH name: ResetSummaryDataSequence :exec
-- DELETE FROM SQLITE_SEQUENCE WHERE name='table_name';

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
