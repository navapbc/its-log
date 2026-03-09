-- params: key_id
WITH total_count (TC) AS (
    SELECT count(*) from itslog_events
)
INSERT INTO itslog_summary
    (key_id, operation, tags, count)
VALUES
    (?, 'count.total', NULL, (select TC from total_count))
