CREATE TABLE IF NOT EXISTS itslog (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    -- automatically provided by the SQLite engine
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- event tag
    event       TEXT NOT NULL,
    -- event value
    -- see https://stackoverflow.com/a/53119060
    value       TEXT NOT NULL,
    -- the type the user claimed for the value
    type        INTEGER NOT NULL
);