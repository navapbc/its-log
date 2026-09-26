WITH
counts AS (
    SELECT key_id, :date as date, 'count.by_tags' as operation, tags, '', count(*) as count
    FROM itslog_events
    GROUP BY tags
)
INSERT OR REPLACE INTO itslog_summary
    (key_id, date, operation, tags, value, count)
SELECT * FROM counts
-- Everything must have gone fine. 
-- Return 0 for success.
RETURNING 0;
