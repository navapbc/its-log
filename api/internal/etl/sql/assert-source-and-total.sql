-- We should get the same count from the total as summing all of the
-- events broken out over sources.
-- Because SQLite stores all numbers as REALs, we need to subtract, take the ABS, 
-- and then ask if the value is less than 0.1... see 
-- https://www.cl.cam.ac.uk/teaching/1213/FPComp/fpcomp12slides.pdf
-- for a reminder of the dangers of floating point.
-- THEREFORE, we return 0 (success) if our delta is less than 0.1, accounting
-- for any REAL number oddities.
SELECT
	(ABS(
		(select count as total from itslog_summary
			where operation = 'count.total' and date = :date)
		-
		(select sum(count) as s from itslog_summary
			where operation = 'count.by_tags' and date = :date)) > 0.1);