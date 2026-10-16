INSERT INTO itslog_summary (
    key_id, date, operation, tags, value, count
    ) VALUES (
    $key_id, $date, $operation, $tags, $value, $count
    )
RETURNING 0;