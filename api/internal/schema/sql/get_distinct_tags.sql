-- name: GetDistinctTags :many
SELECT DISTINCT tags
FROM itslog_events;