DELETE FROM itslog_summary   WHERE 
    key_id = 'ITSLOG_KEY_ID' 
    AND date = 'ITSLOG_DATE'
    AND operation = 'count.total';
