CREATE TABLE IF NOT EXISTS itslog_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- some apps have multiple internal sources
    source      INTEGER NOT NULL,
    -- event tag
    event       INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS itslog_dictionary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- If we reboot in the middle of the day, it is possible
    -- that we will get new hashes for event names. Therefore,
    -- we could get the same event name mapped to multiple hash values.
    -- This is fine. We have to make sure the hash values are unique,
    -- not the names.
    source_name TEXT NOT NULL,
    event_name TEXT NOT NULL,
    source_hash INTEGER NOT NULL,
    event_hash INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON itslog_dictionary (source_hash, event_hash);

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation TEXT NOT NULL,
    source TEXT NOT NULL,
    event TEXT,
    value REAL NOT NULL
);

-- This lets us write and rewrite the data to the same table. 
-- The operation (e.g. `app.by-day`) source (`app_001`) and event (`v2_api`)
-- are "unique", and we then want to record/update the value.
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (operation, source, event);

CREATE TABLE IF NOT EXISTS itslog_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL,
    value TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS metadata_ndx ON itslog_metadata (key);
