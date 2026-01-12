-- events by user
DELETE FROM itslog_summary WHERE operation = 'distinct.by_day.values';

WITH dc AS (
  SELECT count(DISTINCT value_hash) AS cnt FROM itslog_events
)
INSERT INTO itslog_summary 
  (key_id, date, operation, source_name, event_name, value)
SELECT
  'ITSLOG_KEY_ID' as key_id, 'ITSLOG_DATE' as date, 'distinct.by_day.values', NULL as source_name, NULL as event_name, dc.cnt as value
FROM dc
;