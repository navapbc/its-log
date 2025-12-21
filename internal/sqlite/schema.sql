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

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation TEXT NOT NULL,
    source TEXT NOT NULL,
    event TEXT,
    value REAL NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON itslog_dictionary (source_hash, event_hash);
