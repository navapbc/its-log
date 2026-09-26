-- name: GetAllSummaries :many
SELECT * FROM itslog_summary;

-- name: InsertSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  key_id, date, operation, tags, value, count, hash
  ) VALUES (
  ?, ?, ?, ?, ?, ?, ?
  );
