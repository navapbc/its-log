
-- -- Remove previous values for this computation
-- DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_source';

-- -- Compute the counts per event source
-- WITH 
-- counts AS (
--   SELECT 'count.by_day.by_source' as operation, ie.source_hash, count(*) as event_count
--   FROM itslog_events ie
--   GROUP BY ie.source_hash),
-- distinct_names AS (
--   SELECT distinct(id.source_hash), id.source_name 
--   FROM itslog_dictionary id),
-- final AS (
--     SELECT c.operation, dn.source_name, c.event_count
--     FROM counts c
--     JOIN distinct_names dn
--     WHERE dn.source_hash = c.source_hash
--   )
-- INSERT INTO itslog_summary 
--     (operation, source_name, event_name, value)
-- SELECT 
--     operation, source_name, NULL, event_count 
-- FROM final;

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
    SELECT d.source_name, d.event_name, c.event_count
    FROM counts c
    INNER JOIN itslog_dictionary AS d ON d.source_hash = c.source_hash
	  GROUP BY d.source_name
  )
INSERT INTO itslog_summary 
    (operation, source_name, event_name, value)
SELECT 
    'count.by_day.by_source', source_name, NULL, event_count 
FROM final;


