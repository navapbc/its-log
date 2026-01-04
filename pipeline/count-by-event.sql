-- Remove previous values for this computation
DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_source.by_event';

-- Compute the counts per event source
WITH 
counts AS (
  SELECT ie.source_hash, ie.event_hash, count(*) as event_count
  FROM itslog_events ie
  GROUP BY ie.source_hash, ie.event_hash),
distinct_source_names AS (
  SELECT distinct(id.source_hash), id.source_name 
  FROM itslog_dictionary id),
distinct_event_names AS (
  SELECT distinct(id.event_hash), id.event_name
  FROM itslog_dictionary id),
final AS (
    SELECT 'count.by_day.by_source.by_event' as operation, dsn.source_name, den.event_name, c.event_count
    FROM counts c
    JOIN distinct_source_names dsn, distinct_event_names den
    WHERE dsn.source_hash = c.source_hash AND den.event_hash = c.event_hash
  )
INSERT INTO itslog_summary 
    (operation, source_name, event_name, value)
SELECT 
    operation, source_name, event_name, event_count 
FROM final;