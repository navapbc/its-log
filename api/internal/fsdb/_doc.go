// Package fsdb provides an in-filesystem database for its-log
//
// This package does a lot of lifting.
//
// The database backend is SQLite. This gives us 1-2 orders of magnitude performance over
// using a networked database (e.g. Postgres), and allows us to easily backup/etc. what we log.
// It places constraints on us, though: the database can handle near endless concurrent reads,
// but does very poorly with concurrent writes. For this reason, its-log must serialze all writes
// to the database when logging events.
//
// When we create a new database file (which we potentially do with any given event), we create
// the tables and load default ETL actions into the ETL table. These default actions are automatically
// used on all tables for summarization; end-users can load additional actions as desired.
//
// All (most?) database actions take place through sqlc. This package parses the SQL for the table
// schemas and queries, and then turns that SQL into typed Golang functions. So, it is *not* an ORM.
// We write the SQL, and that is then used to generate type-safe Golang functions to interact with the
// database. Therefore, `schema.sql` and `query.sql` are critical to our work with the SQLite files.
package fsdb
