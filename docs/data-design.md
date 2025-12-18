# data design

The `its-log` data design is opinionated and intentional.

## design goals

The design is imagined against a backdrop of compliance-burdened systems that are run on limited budgets. 

* **High throughput**. Ideally, a single `its-log` instance serves many applications. It should be possible for a single `its-log` instance to handle tens of thousands of events *per second*.
* **Low resource requirements**. A single instance should be able to run comfortably on a small host; at most, several hundred megabytes of RAM and megabytes of local disk.
* **No real-time**. `its-log` operates around the assumption that no data visibility is needed regarding *today*. Yesterday is good enough.
* **Semantic compression**. The data collected daily should be "compressed" daily based on knowledge of the data. Because we are counting *events*, it should be possible to compress (say) 100K events of a given type to a single row representing the summation over the day, or (at most) 24 rows, representing an hour-by-hour summation. There should be no need for storing the original, second-by-second raw data.
* **Zero PHI/PII**. `its-log` provides support for encrypting data "one way," guaranteeing that data is countable while eliminating the information content of the event.

Taken together, this leads to a highly performant event-logging system that is intended, primarily, for operational awareness. "How much traffic did each API endpoint in our system receive?" "Are there days we receive more total traffic than others?" "How many unique users 

## api design

The API itself currently has one endpoint; a second might be added. This is described in [api.md](api.md).

## design implementation

`its-log` has two tables under the hood.

### events

The events table stores four values:

* id: the SQLite-native row ID; primary key, auto-incrementing, auto-assigned
* timestamp: the native `DATETIME` type; auto-assigned
* source: a 64-bit signed integer
* event: a 64-bit signed integer

The `its-log` API server consumes a string for the source and event; it then hashes that value, and takes a signed 64-bit portion to write to the database. This means that each row of data becomes roughly 8 + 6 + 8 + 8 bytes, or 30 bytes. These values are then buffered, so that writes happen in bulk; this is configurable, but the default configuration buffers 2000 events and flushes the cache once per minute. (2000x30 bytes is 60K + fixed overhead, meaning we're not taxing RAM to buffer that many events.) The events are then written to disk in a single transaction, again a significant performance optimization.

Alone, this table is indescipherable. It is possible to tell how many unique applications there are, and how many events per app... but, they are meaning-free integers.

### dictionary

The dictionary table is much smaller. It is maintained to map the space-efficient integer representation back to human-readable text.

* id: a row ID/PK, not used
* source_hash: the hashed value of an event source, the same as would be stored in the `source` column of the `events` table
* event_hash: the hashed value of an event, the same as would be stored in the `event` column of the `events` table
* source_name: the plain-text value for the source name associated with this hash
* event_name: the plain-text value for the event name associated with this hash

This dictionary is necessary if events are going to be converted back to human-readable form. However, it is a small table: 30 applications with 100 unique events each would yield a table with 30x100 rows, or 3000 rows. This is tiny, and easily `JOIN`ed with the `events` table as needed for dynamic queries, or used for semantic compression on a daily basis (where the plain text can be used, because there are so many fewer rows).

A new dictionary is started every day. This allows for hash seeds to change over time, and as a result, we maintain a complete mapping from the hashed values to something human-understandable.

## summary

`its-log` stores events. It is intended to support counting things. Those might be events that happen repeatedly (e.g. "how many times was this web page accessed?"), or they might be counts of unique events within the system ("how many unique applications were used each day of the last week?"). It is a spaec-efficient design, capable of storing millions of events in a handful of megabytes of space, and compressing data *meaningfully* down to a handful of rows on a daily basis. 