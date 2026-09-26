-- name: GetETL :one
SELECT name, kind, body, last_run
FROM itslog_etl
WHERE
  name = ?
LIMIT 1;