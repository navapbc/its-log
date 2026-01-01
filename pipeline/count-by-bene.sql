-- events by user
DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_bene';

-- Compute the counts per bene
-- sources look like: blue.fhir.v3
-- events look like: Coverage.AaNeel - CCA.76e9493b
WITH 
  benes AS (
    SELECT distinct event_hash, SUBSTR(SUBSTR(event_name, INSTR(event_name, '.') + 1), INSTR(SUBSTR(event_name, INSTR(event_name, '.') + 1), '.') + 1) AS bene
    FROM itslog_dictionary
  ),
  count_per_bene AS (
    select benes.bene, count(ie.event) as event_count
    from itslog_events ie
    join benes
    where 
      benes.event_hash = ie.event
    group by benes.bene, ie.event
  ),
  final as (
    select bene, sum(event_count) as bene_count from count_per_bene group by bene
  )
INSERT INTO itslog_summary 
    (operation, source, event, value)
SELECT 
    'count.by_day.by_bene', bene, NULL, bene_count 
FROM final;