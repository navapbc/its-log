# data design

The `its-log` data design is opinionated and intentional.

## design goals

The design is imagined against a backdrop of compliance-burdened systems that are run on limited budgets. 

* **High throughput**. Ideally, a single `its-log` instance serves many applications. It should be possible for a single `its-log` instance to handle hundreds of events per second, or tens of thousands of events per minute.
* **Low resource requirements**. A single instance should be able to run comfortably on a small host; at most, several hundred megabytes of RAM and low gigabytes of local disk.
* **No real-time**. `its-log` operates around the assumption that no data visibility is needed regarding *today*. Yesterday is good enough.

Taken together, this leads to a highly performant event-logging system that is intended, primarily, for product- or application-level awareness. "How much traffic did each API endpoint in our system receive?" "Are there days we receive more total traffic than others?" "How many unique users 

## api design

The logging API is essentially a single endpoint. 

There are additional endpoints for loading ETL transforms (SQL) as well as running individual ETL steps or sequences of steps.

## design implementation

`its-log` has several tables under the hood. Two tables drive the data collection, and one drives the ETL.

### events

The events table stores four values:

* id: the SQLite-native row ID; primary key, auto-incrementing, auto-assigned
* timestamp: the native `DATETIME` type; auto-assigned
* key_id: the API key logging the event
* cluster: a string for a unique id connecting *several* events; possibly a UUID or similar
* tags: a string (period separated) representing one or more categorical tags associated with the event (e.g. ["v3", "endpoint"] gets alphabetized and mapped to "endpoint.v3")
* value: a string representing a unique value associated with this particular event

## summary

`its-log` stores events. It is intended to support counting things. Those might be events that happen repeatedly (e.g. "how many times was this web page accessed?"), or they might be counts of unique events within the system ("how many unique applications were used each day of the last week?"). It is a spaec-efficient design, capable of storing millions of events in a handful of megabytes of space, and compressing data *meaningfully* down to a handful of rows on a daily basis. 

The summary table stores:

* last_run: a timestamp for when the row in question was last generated
* key_id: the API generating the summary
* operation: the ETL action that led to this row (e.g. `count.total` vs. `count.by-tag`)
* value: for operations that count across unique values
* count: the number of rows or events counted by this operation

## ETL table

The ETL table stores the ETL events used for transforming the data. The events are recorded in the table for multiple reasons:

1. Historical awareness. It is possible to see the source for each ETL transform that was applied at this point in history. So, if an ETL action changes over time, we can see which version was used at any given point (if necessary).
2. Automation. `its-log` provides several "common" or "standard" ETL actions, but individual applications can load additional ETL actions via API for use on their data.  

The table contains:

* inserted: when the ETL action was inserted into the table
* last_run: when the ETL action was last executed against the data
* key_id: the API key triggering the action (or `its-log` for default actions)
* name: the name of the ETL action
* kind: whether it is a native `sql` action, a `golang` action, or a `sequence` of actions
* body: the body of the action

`sql` actions are typically SQL transforms that read from the `events` table and write to the `summary` table. They should always end with `SELECT 1;` to indicate success. Because all actions are run in a transaction, they will either complete fully or fail. More guidelines on writing ETL actions can be found in [etl.md](etl.md).

`golang` actions are compiled into `its-log`. They carry out transforms that would either be too complex in SQL, or play a role in moving data around. For example, an ETL action that copies the database to S3, or exports a copy of the database as a CSV would be an ETL action that must be expressed in Golang.

`sequence` is a comma- or newline-separated list of sequence actions. Each action should be a valid name in the table. When triggered as a sequence, each ETL action is run in turn; if one fails, the sequence exits. Each individual action is run in a transaction, but not the sequence as a whole.