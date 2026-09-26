WITH
counts AS (
    SELECT key_id, $date as date, 'count.by_tags' as operation, tags, '', count(*) as count
    FROM itslog_events
    GROUP BY key_id, tags
)
INSERT INTO itslog_summary
    (key_id, date, operation, tags, value, count)
SELECT * FROM counts
ON CONFLICT (date, operation, tags, value)
DO UPDATE SET 
    date = EXCLUDED.date,
    operation = EXCLUDED.operation,
    tags = EXCLUDED.tags,
    value = EXCLUDED.value,
    count = EXCLUDED.count
-- Everything must have gone fine. 
-- Return 0 for success.
RETURNING 0;
