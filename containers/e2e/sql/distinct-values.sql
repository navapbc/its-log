-- events by user
DELETE FROM itslog_summary WHERE operation = 'distinct.by_day.values';

WITH dc AS (
  SELECT count(DISTINCT value_hash) AS cnt FROM itslog_events
)
INSERT INTO itslog_summary 
  (key_id, operation, source_name, event_name, value)
SELECT
  ? as key_id, 'distinct.by_day.values', NULL, NULL, dc.cnt
FROM dc
;