# etl

When we're done with a day's worth of data, we want to transform it from it's raw form (many events) to a useful form (condensed, calculated values).

In databases, this is called an "ETL pipeline," which is an acronym for "extract, transform, and load." 

* We extract data from one source (our event tables)
* We transform it (by counting events)
* We load it (into a new table)

## etl in its-log

`its-log` has a light pipeline ETL approach that is experimental at this time. It allows a sequence of events to be defined, and those events can be run by `its-log` using the `etl` subcommand

### jsonnet -> json

We start with a Jsonnet file:

```jsonnet
{
  actions: [
    {
      action: "message",
      message: "helo",
    },
    {
      action: "sql",
      filename: "count-by-source.sql",
    },
    {
      action: "fileCopy",
      source: "env.sourcePath",
      destination: "env.destPath",
    },
  ],
}
```

This is transformed into a JSON document of the same structure.

### its-log etl

We then consume this runscript and an SQLite database.

```
its-log etl --runscript pipeline.json --sqlite 2025-12-24.sqlite
```

This runs each action one-by-one.

## etl actions

There is a small vocabulary of ETL actions.

### message

This emits a log message. Useful for logging/reporting.

### sql

The SQL action runs a standalone SQL script. An example script might be:

```sql
-- Remove previous values for this computation
DELETE FROM itslog_summary WHERE operation = 'count.by_day.by_source';

-- Compute the counts per event source
WITH 
counts AS (
  SELECT 'count.by_day.by_source' as operation, source, count(*) as event_count
  FROM itslog_events
  GROUP BY source),
distinct_names AS (
  SELECT distinct(source_hash), source_name 
  FROM itslog_dictionary),
final AS (
    SELECT operation, source_name, event_count
    FROM counts
    JOIN distinct_names
    WHERE distinct_names.source_hash = counts.source
  )
INSERT INTO itslog_summary 
    (operation, source, event, value)
SELECT 
    operation, source_name, NULL, event_count 
FROM final;
```

This script removes a set of values from the summary table, and then computes a count on a source-by-source basis. It stores the result of those counts into the summary table.

### not yet implemented

Not yet implemented are multiple possible actions:

* fileCopy
* tableToDB
* fileToS3

which would allow us to (say) copy a file, copy a table to a remote database, or move a file to S3. All of these would likely interact with the top-level configuration file, allowing safe interaction with services in production environments.

## open questions

There are at least three ways to tackle this:

* Write the ETL as SQL scripts (as demonstrated here)
* Write the ETL in Go, and incorporate it directly into `its-log`
* Write the ETL in Python, and assume it is executed separately

The benefit to a light ETL pipeline processor in `its-log` is that it will be "in-place" and approved, and therefore extending it to process a JSON file and execute SQL (which can undergo testing and code review) has many benefits.