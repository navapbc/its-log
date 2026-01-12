DELETE FROM itslog_summary WHERE operation = 'distinct.by_day.by_source.values';

WITH counts AS (
  SELECT DISTINCT source_hash, value_hash, count(value_hash) AS the_count 
  FROM itslog_events
  GROUP BY source_hash, value_hash)
INSERT INTO itslog_summary 
  (key_id, date, operation, source_name, event_name, value)
SELECT 
  'ITSLOG_KEY_ID' as key_id, 'ITSLOG_DATE' as date, 'distinct.by_day.by_source.values' as operation, d.source_name as source_name, look.name as event_name, counts.the_count as value
  FROM counts
  JOIN
    itslog_dictionary d, itslog_lookup look
  WHERE 
    counts.source_hash = d.source_hash
  AND
    counts.value_hash = look.hash
  GROUP BY look.name
;