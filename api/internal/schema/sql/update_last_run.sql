-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = unixepoch() 
WHERE 
  name = ?;-- name: UpdateLastRun :exec
UPDATE itslog_etl
  SET 
    last_run = unixepoch() 
WHERE 
  name = ?;