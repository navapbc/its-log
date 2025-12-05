# its-log

*It's better than bad, it's good!*

Data collection and analysis platforms are a kind of lock-in: once you buy into them, and put all your data in, getting out is a complex and expensive proposition.

## what is `its-log` best at?

`its-log` is best suited to low-to-medium volume event logging. It is *not*:

* a replacement for ElasticSearch or other high-volume systems
* intended to handle complex application events
* suitable for real-time analysis

It **is** well suited to:

* capturing application events
* storing categorical/tagged numerical data
* exporting data for daily/weekly batch processing

`its-log` uses open data formats, and is intended to be the start of small, purpose/product-driven analytics pipelines.

## how do I log to `its-log`

Once built, `its-log` is a single binary executable that provides an API with a single endpoint:

```
POST /<version>/log
```

The endpoint expects a JSON structure that MUST match this shape:

```json
{
    "event": string,
    "value": string,
    "type": integer
}
```

Valid values for `type` are:

| type | meaning |
| -- | -- |
| 1 | INTEGER |
| 2 | REAL |
| 3 | TEXT | 
| 4 | DATE |
| 5 | DATETIME |
| 6 | JSONB |
| 7 | BLOB | 

These types are to help application developers decode the data when extracting it from SQLite/JSON, nothing more.

All events are timestamped internally by `its-log`.

## how do I get data out?

`its-log` stores its data to one of two places:

1. An SQLite database
1. An S3-compatible bucket

Each produces similar-but-different data.

### what does the SQLite data look like?

SQLite is an in-filesystem database. What this means is that data is stored into just one file. Because there are practical (and real) limits to how large we are allowed to make one file, `its-log` creates a new file every day. Therefore, in order to get the data from Monday, December 1st, we need to get the file `2025-12-01-v1.itslog`.

It's log can either leave these files on disk (to be backed up by some other process), or it can automatically copy these files into an S3 bucket every day, deleting the local file when the copy completes successfully.

### what does the S3 data look like?

When storing data to an S3-compatible environment, `its-log` stores every single event separately as a JSON document. It keys them by the day, so that the bucket containing the data looks like:

```
2025-12-01/
    |- 0000000000.json
    |- 0000000001.json
    |- ...
    |- 8640000000.json
2025-12-02/
    |- 0000000000.json
    |- ...
```

where each file is named with a simple counter. `its-log` assums you will never capture more than 100,000 events in a single second. Which is good, because it can't.

## how do I process `its-log` data?

If working with SQLite databases, it is possible to download these files and browse them (using any number of free/open tools), or code can be written to process the data using Python. A typical worklow might:

1. Download an SQLite file
1. Copy its contents into a Postgres database
1. Analyze/process the data from there

If working with S3 buckets, the workflow might look like:

1. Read each JSON file out of S3 from a given day
1. Load the data entry into Postgres
1. Analyze/process the data from there

## this is <insert negative word>! we should just use <expensive proprietary thing>!

Yes. You could. You probably already have. And now you have problems 

1. Accessing your data
1. Working with your data
1. Migrating your data

and so on. This is meant to be the simplest possible starting point for your application log data. The purpose of `its-log` is to avoid lock-in by starting with the (nearly) simplest possible data capture strategy.

## hiow do I deploy `its-log`?

It is probably easiest to put the executable into a container, run it, and expose its port. If the container has an immutable filesystem, then the only option for data storage will be S3.

## how do I configure `its-log`?

Everything is configured with a single YAML file. All security-critical variables can be overridden by the environment. This way, local developers can put (safe, insecure) values in for testing, and CI/CD processes or production environments can override those values during testing and deployment.

## packages used

* [github.com/spf13/viper](https://github.com/spf13/viper)
* go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
* go get gocloud.dev
