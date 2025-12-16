# 1. Packages used for development

Date: 2025-12-05

## Status

Accepted

## Context

Choices needed to be made for the base API and storage layer.

## Decision

* The logger will be written in Go; this provides compiled performance, typing, and CSP-based concurrency primitives (channels, etc.)
* `gin` will be used as the API framework, being widely used and well documented
* `slqc` will be used to manage the database, as it allows for pure-SQL expression of tables and queries that are then "lifted" or "compiled" into the host language, with typing
* `modernc` is a pure golang SQLite driver, and will be used instead of a driver with binary dependencies (e.g. `mattn`), to minimize architectural impacts if a move between Intel and ARM platforms is needed
* `sqlite` will be used for the local database technology, as it provides excellent reliability and throughput

## Consequences

None. Or, all choices have impact. Library choices can change if needed.