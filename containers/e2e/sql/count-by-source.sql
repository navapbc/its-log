-- Remove previous values for this computation
DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_source';

-- Compute the counts per event source
WITH 
counts AS (
  SELECT ie.source_hash, ie.event_hash, count(*) as event_count
  FROM itslog_events ie
  GROUP BY ie.source_hash
  ),
final AS (
    SELECT 'ITSLOG_KEY_ID' as key_id, 'ITSLOG_DATE' as date, 'count.by_day.by_source' as operation, 
        d.source_name as source_name, NULL as event_name, c.event_count as value
    FROM counts c
    INNER JOIN itslog_dictionary AS d ON d.source_hash = c.source_hash
	  GROUP BY d.source_name
  )
INSERT INTO itslog_summary 
    (key_id, date, operation, source_name, event_name, value)
SELECT * FROM final;


