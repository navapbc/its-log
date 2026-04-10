# etl

When we're done with a day's worth of data, we want to transform it from it's raw form (many events) to a useful form (condensed, calculated values).

In databases, this is called an "ETL pipeline," which is an acronym for "extract, transform, and load." 

* We extract data from one source (our event tables)
* We transform it (by counting events)
* We load it (into a new table)

## etl in its-log

`its-log` has a light pipeline ETL approach that is experimental at this time. It allows a sequence of events to be defined, and those events can be run by `its-log` using the `etl` subcommand

### a sequence

`its-log` supports loading a sequence into the ETL table. This might look like

```
sentinel
count-total
count-by-tags
count-all-combinations
assert-source-and-total
```

where each line is the name of an ETL action in the table. This can also be a single line with commas.

### sql ETL actions

All SQL ETL actions:

* Are run under a transaction; they either complete fully or fail
* Must return 0 (`false`, or failure) or 1 (`true`, or success)
* All ETLs are passed the `:key_id`, the `:app_id`, and the `:date` as named parameters
* The summary table has a uniqueness criterion spanning (`date`, `operation`, `tags`, and `value`)

Therefore, any action that reads from the events table and writes to the summary table should end with a single statement of `SELECT 1;`. If the transform succeeds, this has the effect of returning a success value from the transaction as a whole.

Note that each application has its own database; it is not possible for the ETL actions of one application to operate, even accidentally, on the data of other applications.

#### defaults: sentinel 

The `sentinel` action is a no-op, is provided as a default, and is used by `its-log` to verify the existence of the table.

It looks like this:

```
SELECT 1;
```

Running this as a first step in any ETL action is good practice.

#### defaults: count-total

```
WITH total_count (TC) AS (
    SELECT count(*) from itslog_events
)
INSERT OR REPLACE INTO itslog_summary
    (key_id, operation, tags, count)
VALUES
    (:key_id, 'count.total', NULL, (select TC from total_count));
SELECT 1;
```

This demonstrates several things.

* A CTE is used to perform the initial select.
* The data is inserted into the summary table as `INSERT OR REPLACE`. More on this below.
* The `:key_id` is used to denote which API triggered this action.
* It returns 1 to indicate the transaction succeeded.

We use `INSERT OR REPLACE` to make sure there is only one `count.total` in the summary table. This works in part because of the following index on the summary table:

```
CREATE UNIQUE INDEX IF NOT EXISTS summary_ndx ON itslog_summary (date, operation, IFNULL(tags, 0), IFNULL(value, 0));
```

Because the events table only contains events from a single day (or, it should by design), it is sufficient to count all of the events in the table and consider them a total count for the day.

#### defaults: count-by-tags

```
WITH
counts AS (
    SELECT 'count.by_tags' as operation, tags, count(*) as count
    FROM itslog_events
    GROUP BY tags
)
INSERT OR REPLACE INTO itslog_summary
    (operation, tags, count, key_id)
SELECT *, :key_id as key_id FROM counts;
SELECT 1;
```

This ETL step gathers the unique tags in the table and inserts a count for each key into the summary, demonstrating a `GROUP BY` ETL step.

#### defaults: count-all-combinations

This is an example of a golang-implemented ETL step. The source is included in the database, but it is actually executed from the compiled code. (There is no way to dynamically insert a golang-implemented ETL step.)

This step 

1. Gathers all of the unique tags in the events table
2. Splits them all
3. Builds unique combinations of all tags
4. Counts the presence of those combinations
5. Inserts all non-zero counts into the summary table.

This would be difficult to do in SQL, but combinations-generation is easy in a general purpose language. 

## future ETL steps

Examples of additional ETL steps might include:

* backup-to-s3: a golang-implemented ETL step could easily be written such that it would safely backup the current database to S3. 
* export-to-csv: a golang-implemented ETL that could take the events and/or summary table and export it to a CSV file

There are others; steps written in Golang necessarily have access to the full environment, and therefore can interact with other systems, etc. 

There is no way to pass parameters to ETL steps. This guarantees some level of safey and consistency.
