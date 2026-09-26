-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = epoch(now()) 
WHERE 
  name = ?;
