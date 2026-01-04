-- We should get the same count from the sum over the events by source
-- and the events split over all of sub-categories of events.
-- Returns 1 if true, 0 if false.
SELECT
	(select sum(value) as s from itslog_summary
		where operation = 'count.by_day.by_source')
	=
	(select sum(value) as s from itslog_summary
		where operation = 'count.by_day.by_source.by_event') as eq
