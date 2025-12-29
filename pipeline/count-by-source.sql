
-- Remove previous values for this computation
DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_source';

-- Compute the counts per event source
WITH 
counts AS (
  SELECT 'count.by_day.by_source' as operation, source, count(*) as event_count
  FROM itslog_events
  GROUP BY source),
distinct_names AS (
  SELECT distinct(source_hash), source_name 
  FROM itslog_dictionary),
final AS (
    SELECT operation, source_name, event_count
    FROM counts
    JOIN distinct_names
    WHERE distinct_names.source_hash = counts.source
  )
INSERT INTO itslog_summary 
    (operation, source, event, value)
SELECT 
    operation, source_name, NULL, event_count 
FROM final;

-- This could be part of some kind of check or assertion?
-- We would need to know what we expect from our data.
-- select * from itslog_summary where operation = 'count.by_day.by_source' order by source;