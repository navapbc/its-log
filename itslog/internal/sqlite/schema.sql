CREATE TABLE IF NOT EXISTS itslog_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- cluster is useful for a related set of events
    cluster_hash INTEGER,
    -- some apps have multiple internal sources
    source_hash INTEGER NOT NULL,
    -- event tag
    event_hash  INTEGER NOT NULL,
    -- value is useful for when you want a unique value 
    -- associated with this event
    value_hash INTEGER
);

CREATE TABLE IF NOT EXISTS itslog_dictionary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    source_hash INTEGER NOT NULL,
    source_name TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_hash INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON itslog_dictionary (source_hash, event_hash);

CREATE TABLE IF NOT EXISTS itslog_lookup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    hash INTEGER NOT NULL,
    name TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS lookup_hashes_ndx ON itslog_lookup (hash);

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    date DATE DEFAULT CURRENT_DATE NOT NULL,
    operation TEXT NOT NULL,
    source_name TEXT,
    event_name TEXT,
    value REAL NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (operation, source_name, event_name);

CREATE TABLE IF NOT EXISTS itslog_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    key INTEGER NOT NULL,
    value TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS key_hashes_ndx ON itslog_metadata (key);

CREATE TABLE IF NOT EXISTS itslog_etl (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key_id TEXT NOT NULL,
    inserted DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    last_run DATETIME,
    sql TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS step_name_hashes_ndx ON itslog_etl (name);


-------------------------------------------------------------------------------------
-- SUMMARIES
-------------------------------------------------------------------------------------

-- -- count.total
-- DROP VIEW IF EXISTS itslog_summary_count_total;
-- CREATE VIEW itslog_summary_count_total AS
-- WITH tot AS
--     (
-- 		SELECT count(*) as event_count from itslog_events
-- 	)
-- SELECT 
--     'count.total' as operation, NULL as source_name, NULL as event_name, tot.event_count
-- FROM tot;

-- -- count.by_day.by_source
-- DROP VIEW IF EXISTS itslog_summary_count_by_source;
-- CREATE VIEW itslog_summary_count_by_source AS
-- WITH 
-- counts AS (
--   SELECT ie.source_hash, ie.event_hash, count(*) as event_count
--   FROM itslog_events ie
--   GROUP BY ie.source_hash
--   ),
-- final AS (
--     SELECT d.source_name, d.event_name, c.event_count
--     FROM counts c
--     INNER JOIN itslog_dictionary AS d ON d.source_hash = c.source_hash
-- 	  GROUP BY d.source_name
--   )
-- SELECT 
--     'count.by_day.by_source' as operation, source_name, NULL as event_name, event_count 
-- FROM final;

-- -- count.by_day.by_source.by_event
-- DROP VIEW IF EXISTS itslog_summary_count_by_event;
-- CREATE VIEW itslog_summary_count_by_event AS
-- WITH 
-- counts AS (
--   SELECT ie.source_hash, ie.event_hash, count(*) as event_count
--   FROM itslog_events ie
--   GROUP BY ie.source_hash, ie.event_hash),
-- final AS (
--     SELECT d.source_name, d.event_name, c.event_count
--     FROM counts c
--     INNER JOIN itslog_dictionary AS d ON d.source_hash = c.source_hash AND d.event_hash = c.event_hash
-- 	GROUP BY d.source_name, d.event_name
--   )
-- SELECT 
--     'count.by_day.by_source.by_event' as operation, source_name, event_name, event_count 
-- FROM final;
