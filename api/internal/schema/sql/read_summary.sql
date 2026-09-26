-- name: ReadSummary :one
SELECT 
  date, 
  operation, 
  tags, 
  value,
  count
FROM itslog_summary
WHERE 
  tags LIKE ?
  AND
  operation LIKE ?
ORDER BY id
LIMIT 1
;