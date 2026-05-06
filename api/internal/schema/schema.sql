CREATE TABLE IF NOT EXISTS itslog_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp INTEGER NOT NULL DEFAULT (unixepoch()), 
    -- so we know what key performed the operation
    key_id TEXT NOT NULL,
    -- cluster is useful for a related set of events
    cluster TEXT,
    -- some apps have multiple internal sources
    tags TEXT NOT NULL,
    -- value is useful for when you want a unique value 
    -- associated with this event
    value TEXT
);

CREATE TRIGGER IF NOT EXISTS itslog_events_prevent_delete
BEFORE DELETE ON itslog_events
BEGIN
    SELECT RAISE(ABORT, 'no deletion from events allowed');
END;

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    last_updated INTEGER DEFAULT (unixepoch()) NOT NULL,
    date TEXT NOT NULL,
    key_id TEXT NOT NULL,
    operation TEXT NOT NULL,
    tags TEXT NOT NULL,
    value TEXT NOT NULL,
    count REAL NOT NULL,
    hash TEXT
);

-- See https://stackoverflow.com/questions/22699409/sqlite-null-and-unique
-- We want NULL values to count towards uniqueness here.
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (date, operation, IFNULL(tags, 0), IFNULL(value, 0));

CREATE TABLE IF NOT EXISTS itslog_etl (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inserted INTEGER DEFAULT (unixepoch()) NOT NULL,
    last_run INTEGER,
    key_id TEXT NOT NULL,
    name TEXT NOT NULL,
    kind TEXT NOT NULL,
    body TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS step_name_ndx ON itslog_etl (name);

-- CREATE TABLE IF NOT EXISTS itslog_sequences (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     inserted DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     key_id TEXT NOT NULL,
--     name TEXT NOT NULL,
--     steps TEXT NOT NULL
-- );
-- CREATE UNIQUE INDEX IF NOT EXISTS sequence_name_ndx ON itslog_sequences (name);

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
