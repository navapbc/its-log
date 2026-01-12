DELETE FROM itslog_summary WHERE operation = 'distinct.by_day.by_source.values';

WITH counts AS (
  SELECT DISTINCT source_hash, value_hash, count(value_hash) AS the_count 
  FROM itslog_events
  GROUP BY source_hash, value_hash)
INSERT INTO itslog_summary 
  (key_id, operation, source_name, event_name, value)
SELECT 
  ? as key_id, 'distinct.by_day.by_source.values', d.source_name, look.name, counts.the_count
  FROM counts
  JOIN
    itslog_dictionary d, itslog_lookup look
  WHERE 
    counts.source_hash = d.source_hash
  AND
    counts.value_hash = look.hash
  GROUP BY look.name
;