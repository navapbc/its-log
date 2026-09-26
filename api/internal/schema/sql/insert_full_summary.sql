-- name: InsertFullSummary :exec
INSERT OR REPLACE INTO itslog_summary (
  last_run, key_id, date, operation, tags, value, count, hash
  ) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
  );