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

-- name: TestEventPairExists :one
SELECT EXISTS(
  SELECT 1 
  FROM events 
  WHERE 
    source = ?
    AND
    event = ?
  );


-- name: TestDictionaryPairExists :one
SELECT EXISTS(
  SELECT 1 
  FROM dictionary
  WHERE 
    source_hash = ?
    AND
    event_hash = ?
  );
