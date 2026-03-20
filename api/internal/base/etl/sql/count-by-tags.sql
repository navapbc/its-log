WITH
counts AS (
    SELECT key_id, 'count.by_tags' as operation, tags, count(*) as count
    FROM itslog_events
    GROUP BY tags
)
INSERT OR REPLACE INTO itslog_summary
    (key_id, operation, tags, count)
SELECT * FROM counts;
-- Everything must have gone fine. 
-- Return 1 or 'true' for success.
SELECT 1;