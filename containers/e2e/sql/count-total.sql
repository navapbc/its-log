DELETE FROM itslog_summary WHERE operation = 'count.total';

WITH tot AS
    (
		SELECT count(*) as cnt from itslog_events
	)
INSERT INTO itslog_summary 
    (key_id, operation, source_name, event_name, value)
SELECT 
    ? as key_id, 'count.total' as operation, NULL as source_name, NULL as event_name, tot.cnt as value
FROM tot;
