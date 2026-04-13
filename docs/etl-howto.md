# writing your own ETL steps

So, you want to transform some data. Great! There are three ways.

## what is an ETL step?

First, what is an ETL step?

The purpose of an ETL is to drive upstream vizualizations, dashboards, or reports. Therefore, we want to take a large number of `events`, and produce a small number of `summary` rows.

`its-log` expects ETL steps to do one of two things, mostly.

1. Analyze data in `itslog_events` and insert the resulting analysis into `itslog_summary`.
2. Export data from one or more tables to external systems.

`its-log` has a trigger built into the schema to prevent deletions on the `itslog_events` table. There is an index on the `summary` table that encourages overwriting of rows where critical values are the same, enforcing uniqueness.

## using SQL

Imagine we want to know how many events occured on a given day. (Every `its-log` database represents a single day. A new database is automatically created by `its-log` based on the date of an event.) This typically looks like:

```sql
SELECT count(*) from itslog_events
```

We want this to become an ETL step, which means we want to record that value into the `summary` table. We must insert six values:

* key_id
* date
* operation
* tags
* value
* count

When writing an ETL, three values are available to you, and will be passed into your query as variables if you need them:

* :key_id
* :app_id
* :date

As a start, we might do this:

```sql
INSERT INTO itslog_summary
    (key_id, date, operation, tags, value, count)
VALUES
    (:key_id, :date, 'count.total', '', '', (select count(*) from itslog_events))
;
```

However, we want our ETLs to be able to be re-run. As written, this ETL will run into a uniqueness constraint on the `summary` table, and a second run will fail. Therefore, we want to make one critical change:

```sql
INSERT OR REPLACE INTO itslog_summary
    (key_id, date, operation, tags, value, count)
VALUES
    (:key_id, :date, 'count.total', '', '', (select count(*) from itslog_events))
;
```

Note the `INSERT INTO OR REPLACE`. It is generally going to be the case that *all* ETL actions will want to be `INSERT INTO OR REPLACE`, so that when a second `count.total` for a given `date` with the same operation, tags, value, and count, it will replace as opposed to become a duplicate.

Finally, all ETL actions **MUST** return a value.

```sql
INSERT OR REPLACE INTO itslog_summary
    (key_id, date, operation, tags, value, count)
VALUES
    (:key_id, :date, 'count.total', '', '', (select count(*) from itslog_events))
RETURNING 0
;
```

Following the UNIX convention, we return 0 on success.

### the simplest sql ETL step

`its-log` embeds this ETL under the name `sentinel`. It is used as a sentinel value in the database, and can be looked for/executed to make sure that the ETL table is present, populated, and working.

```sql
SELECT 0;
```

### sql ETL assertions

It is also possible to assert conditions that must be true as an ETL step.

```sql
SELECT
	(ABS(
		(select count as total from itslog_summary
			where operation = 'count.total' and date = :date)
		-
		(select sum(count) as s from itslog_summary
			where operation = 'count.by_tags' and date = :date)) > 0.1);
```

This ETL:

1. Looks up two values in the `summary` table
2. Subtracts them
3. Asserts that they must, when subtracted, equal zero.

Because the `count` column is a `FLOAT`, we actually check that the result of the subtraction is less than 0.1. This is because [floating point numbers are hard to get right](https://www.computer.org/csdl/magazine/cs/2014/04/mcs2014040080/13rRUwhpBTR). 

The purpose of this ETL action is to pass quietly if the previous actions did the right thing, or to fail if they did not. In this way, we can assert the correctness of an ETL pipeline as it is executing.


## writing ETL steps in Starlark

[Starlark](https://starlark-lang.org/) is a scripting language used primarily in Google's [Bazel](https://bazel.build/) build tool. It has several properties that make it ideal for use in ETLs:

1. It is deterministic. It is not possible to write an infinite loop in Starlark, so every Starlark program is guaranteed to terminate.
2. It is hermetic. Starlark programs cannot write to disk, talk to the network, or generally break out of their sandbox.
3. It is Python-like. Starlark looks like and generally behaves like a subset of Python.

All Starlark ETL actions have access to one magic function: `query()`. This function takes a string and returns a list of dictionaries. The string is an SQL statement that will be executed in the context of the current database, and the dictionaries have keys matching the columns of of the table queried; all values come back as strings. (FIXME: currently, this *only* works on the `events` table; `its-log` needs to be updated so that queries can be run against `events` or `summary`.)

The data returned can be analyzed/processed arbitrarily, and it is expected that a list of dictionaries will be returned matching the shape of `summary` rows. That data will be processed and inserted after the Starlark code executes.

```python
def summarize():
    query_rows = query("SELECT * from itslog_events")
    summaries = []
    total_vowels = 0
    total_consonants = 0
    for row in query_rows:
        for l in "aeiou".elems():
            if l in row["tags"]:
                total_vowels += 1
        for l in "bcdfghjklmnpqrstvwxyz".elems():
            if l in row["tags"]:
                total_consonants += 1

    result = [
        {"operation": "total.vowels", "count": total_vowels},
        {"operation": "total.consonants", "count": total_consonants}
    ]
    return result
```

The above example demonstrates a query that returns all of t he rows in `events`, counts the vowels and consonants in the `tags` in those rows, and returns two summary rows.

You can install the `starlark` interpeter and test Starlark programs locally. You will need to mock results from query, and call the function directly, but it is possible. In time, it may be that `summarize()` will take parameters, including a parameter that allows you to know whether it is being called in a test or production context. For now, see examples in the codebase for how a developer can comment/uncomment code to make the ETL actions "testable."

## ETL steps in Golang

It is also possible to write ETL steps in Golang. These can work in the same way as an SQL step, analyzing the `events` table and inserting one or more `summary` rows, or they can manipulate the database and external systems (e.g. converting the DB from SQlite to CSV, or copying the database file to another storage medium like S3).

An ETL written in Go must have the following signature:

```golang
func FunName(etlP *types.RunEtlParams) error
```

The RunEtlParams struct contains the following values:

```golang
type RunEtlParams struct {
	AppId   string
	GinCtx  *gin.Context
	EtlName string
	KeyId   string
	Storage *Storage
	Payload map[string]any
}
```

The most important thing to note is the `Payload` dictionary. When an ETL is called, an arbitrary JSON document can be included in the body of the `POST`. In this way, Golang ETLs can be parameterized. It is up to the ETL author to check the correctness of the JSON, and make sure that the params are expected are present. This payload is also available and checked when a `sequence` of ETL actions is called.

Golang ETLs must be compiled into the application, and should generally be considered when 1) the complexity of the ETL exceeds all other approaches, 2) an external system needs to be used as part of the ETL action, or 3) it is of common enough value that it should be provided to all users of `its-log`.
