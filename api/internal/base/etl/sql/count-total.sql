WITH total_count (TC) AS (
    SELECT count(*) from itslog_events
)
INSERT OR REPLACE INTO itslog_summary
    (key_id, operation, tags, count)
VALUES
    (:key_id, 'count.total', NULL, (select TC from total_count));
-- Everything must have gone fine. 
-- Return 1 or 'true' for success.
SELECT 1;