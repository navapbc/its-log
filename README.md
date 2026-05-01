# its-log

*It's better than bad, it's good!*

`its-log` is a lightweight event logger and ETL environment for resource-constrained, compliance-burdened environments.

On a Mac M4, `its-log` can sustain logging 30K events/second to a local SQLite database.


## testing

`its-log` can test itself. It will run sequential correctness tests followed by parallel stress tests, exercising itself end-to-end and validating database contention correctness.

In one terminal, run the storage stack (a [ministack](https://ministack.org/) S3):

```
make storage
```

and in another, run the E2E:

```
make test
```

## debugging its-log

A VS Code debugging bundle is prepared and part of this repository, with all of the envrionment variables required to run. Press `F5` to launch `its-log` under the debugger, and `COMMAND-F5` to halt the debugger (on a Mac).

## testing with Postman

A Postman bundle is included under the `docs` directory. It includes a set of scripts to test results, and can serve as a light E2E test. 

## running its-log

`its-log` can be compiled for multiple architectures.

To run on a ARM processors (Mac):

```
make itslog-arm
```

To run on Intel:

```
make itslog-amd
```

In either case, you will need to make sure a suite of environment variables are present to run `its-log` natively.

## model

`its-log` is opinionated, and espouses ways of thinking about logging and subsequent ETL/analysis pipelines.

### runtime events

Applications generate events which are logged. Any given event:

1. **could** belong to a *cluster* of logs
2. **must** be *categorically tagged*
3. **could** have an associated *unique* value

Some events are logged in real-time, and others might be logged via periodic (hourly/daily/weekly) analysis of internal databases.

### ETL

An explicit goal of `its-log` is to move analysis "left" in the pipeline, and reduce the number of moving parts (e.g. GitHub Actions, AWS Lambdas, etc.) involved in transforming raw events into end-user facing information and visualization. To this end, `its-log` encapsulates a light, SQLite-based ETL pipeline infrastructure, allowing data to be processed within the database and, in doing so:

1. Simplify processing at the end of the pipeline (e.g. in Metabase, Superset, QuickSuite, or similar.)
2. Keeping analytical code (and its results) with the data for archived compliance


## the API

The logging endpoints are rooted at `/v1`. So, `/log/create` should be read as `/v1/log/create`.

| HTTP | Endpoint         | Permissions | Desc                                                                                             |
| ---- | ---------------- | ----------- | ------------------------------------------------------------------------------------------------ |
| POST | /log/create      | logging     | Log an event with cluster (optional), tags, value (optional)                                     |
| POST | /log/create/date | admin       | Log an event per above plus a *date*, for creating test scenarios by logging events in the past. |


The ETL and analysis endpoints are

| HTTP | Endpoint                  | Permissions | Desc                                                                  |
| ---- | ------------------------- | ----------- | --------------------------------------------------------------------- |
| POST | /etl/create               | admin       | Create an ETL action                                                  |
| POST | /etl/run/:date/:name      | admin       | Run an ETL action. Takes optional params in the JSON body.            |
| POST | /sequence/create          | admin       | Create an ETL sequence                                                |
| POST | /sequence/run/:date/:name | admin       | Run a sequence for a given date. Takes optional params like /etl/run. |

For working with the summary table directly

| HTTP | Endpoint      | Permissions | Desc                                                                      |
| ---- | ------------- | ----------- | ------------------------------------------------------------------------- |
| POST | /summary/read | admin       | Read from the summary table; used for testing the results of ETL actions. |


Administrative endpoints include


| HTTP | Endpoint | Permissions | Desc                             |
| ---- | -------- | ----------- | -------------------------------- |
| GET  | /health  | any         | A standard healthcheck endpoint  |
| GET  | /status  | admin       | Get server stats (RAM, GC, etc.) |


## extending its-log

See [DEVELOPMENT.md](docs/DEVELOMENT.md) in `docs`. 

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=navapbc/its-log&type=date&legend=top-left)](https://www.star-history.com/#navapbc/its-log&type=date&legend=top-left)