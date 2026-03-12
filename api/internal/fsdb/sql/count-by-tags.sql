WITH
counts AS (
    SELECT key_id, 'count.by_tags' as operation, tags, count(*) as count
    FROM itslog_events
	WHERE key_id = :key_id
    GROUP BY tags
)
INSERT INTO itslog_summary
    (key_id, operation, tags, count)
SELECT * FROM counts
;