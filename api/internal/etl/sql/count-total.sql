WITH total_count (TC) AS (
    SELECT count(*) from itslog_events
)
INSERT INTO itslog_summary
    (key_id, date, operation, tags, value, count)
VALUES
    ($key_id, $date, 'count.total', '', '', (select TC from total_count))
ON CONFLICT (date, operation, tags, value)
DO UPDATE SET 
    date = EXCLUDED.date,
    operation = EXCLUDED.operation,
    tags = EXCLUDED.tags,
    value = EXCLUDED.value
-- Everything must have gone fine. 
-- Return 0 for success.
RETURNING 0;