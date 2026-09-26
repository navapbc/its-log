SELECT 
  date, 
  operation, 
  COALESCE(tags, '') as tags, 
  value,
  count
FROM itslog_summary
WHERE 
  tags LIKE COALESCE(?, '%')
  AND
  operation LIKE ?
ORDER BY id
;