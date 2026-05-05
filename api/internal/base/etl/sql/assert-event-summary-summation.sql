SELECT
	(ABS(
		(select count from itslog_summary
			where operation = 'total.summary' and date = :date)
		-
		(select sum(count) as s from itslog_summary
			where operation != 'total.summary' and date = :date)) > 0.1);