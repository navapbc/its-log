WITH final AS
    (
        SELECT 'ITSLOG_KEY_ID' as key_id, 'ITSLOG_DATE' as date, 'count.total' as operation, NULL as source_name, NULL as event_name, count(*) as value
        FROM itslog_events
	)
INSERT INTO itslog_summary 
    (key_id, date, operation, source_name, event_name, value)
SELECT * from final;