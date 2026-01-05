DELETE FROM itslog_summary WHERE operation = 'count.total';

WITH tot AS
    (
		SELECT count(*) as cnt from itslog_events
	)
INSERT INTO itslog_summary 
    (operation, source_name, event_name, value)
SELECT 
    'count.total', NULL, NULL, tot.cnt
FROM tot;
