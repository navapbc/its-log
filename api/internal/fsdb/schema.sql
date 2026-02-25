CREATE TABLE IF NOT EXISTS itslog_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- so we know what key performed the operation
    key_id INTEGER NOT NULL,
    -- cluster is useful for a related set of events
    cluster_hash INTEGER,
    -- some apps have multiple internal sources
    tags_hash INTEGER NOT NULL,
    -- value is useful for when you want a unique value 
    -- associated with this event
    value_hash INTEGER
);

CREATE TRIGGER IF NOT EXISTS itslog_events_prevent_delete
BEFORE DELETE ON itslog_events
BEGIN
    SELECT RAISE(ABORT, 'no deletion from events allowed');
END;

-- For mapping hashes back to strings
CREATE TABLE IF NOT EXISTS itslog_lookup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    key_id INTEGER NOT NULL,
    hash INTEGER NOT NULL,
    name TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS lookup_hashes_ndx ON itslog_lookup (hash);

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT DEFAULT CURRENT_DATE NOT NULL,
    key_id INTEGER NOT NULL,
    operation TEXT NOT NULL,
    source_name TEXT,
    event_name TEXT,
    value REAL NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (date, operation, source_name, event_name);

CREATE TABLE IF NOT EXISTS itslog_etl (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inserted DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    key_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    last_run DATETIME,
    sql TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS step_name_hashes_ndx ON itslog_etl (name);


-- CREATE TABLE IF NOT EXISTS itslog_dictionary (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     -- key_id TEXT NOT NULL,
--     timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     source_hash INTEGER NOT NULL,
--     source_name TEXT NOT NULL,
--     event_name TEXT NOT NULL,
--     event_hash INTEGER NOT NULL
-- );
-- CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON itslog_dictionary (source_hash, event_hash);

-- CREATE TABLE IF NOT EXISTS itslog_metadata (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     key_id TEXT NOT NULL,
--     key INTEGER NOT NULL,
--     value TEXT NOT NULL
-- );
-- CREATE UNIQUE INDEX IF NOT EXISTS key_hashes_ndx ON itslog_metadata (key);
