INSERT INTO itslog_summary 
(key_id, date, operation, tags, value, count)
VALUES
(?, ?, 'count.combinations', ?, '', ?)
ON CONFLICT (date, operation, tags, value)
DO UPDATE SET
    date = EXCLUDED.date,
    operation = EXCLUDED.operation,
    tags = EXCLUDED.tags,
    value = EXCLUDED.value
RETURNING 0;