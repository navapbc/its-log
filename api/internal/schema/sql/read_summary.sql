-- name: ReadSummary :one
SELECT 
  date, 
  operation, 
  tags, 
  value,
  count
FROM itslog_summary
WHERE 
  tags LIKE $tags
  AND
  operation LIKE $operation
ORDER BY id
LIMIT 1
;