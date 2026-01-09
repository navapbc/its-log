CREATE TABLE IF NOT EXISTS itslog_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
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
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    source_hash INTEGER NOT NULL,
    source_name TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_hash INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS dictionary_pairs_ndx ON itslog_dictionary (source_hash, event_hash);

CREATE TABLE IF NOT EXISTS itslog_lookup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    hash INTEGER NOT NULL,
    name TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS lookup_hashes_ndx ON itslog_lookup (hash);

CREATE TABLE IF NOT EXISTS itslog_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE DEFAULT CURRENT_DATE NOT NULL,
    operation TEXT NOT NULL,
    source_name TEXT,
    event_name TEXT,
    value REAL NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (operation, source_name, event_name);
