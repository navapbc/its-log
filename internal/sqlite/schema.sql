CREATE TABLE IF NOT EXISTS events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- some apps have multiple internal sources
    source      INTEGER NOT NULL,
    -- event tag
    event       INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS dictionary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- If we reboot in the middle of the day, it is possible
    -- that we will get new hashes for event names. Therefore,
    -- we could get the same event name mapped to multiple hash values.
    -- This is fine. We have to make sure the hash values are unique,
    -- not the names.
    event_source TEXT NOT NULL,
    event_name TEXT NOT NULL,
    source_hash INTEGER NOT NULL,
    event_hash INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON dictionary (source_hash, event_hash);
