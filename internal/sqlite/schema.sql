CREATE TABLE IF NOT EXISTS itslog (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- some apps have multiple internal sources
    source      TEXT NOT NULL,
    -- event tag
    event       TEXT NOT NULL
);