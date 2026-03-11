# its-log

*It's better than bad, it's good!*

`its-log` is a lightweight event logger and ETL environment for resource-constrained, compliance-burdened environments.

On a Mac M4, `its-log` can sustain logging 30K events/second to a local SQLite database.

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

### infastructure as code

`its-log` can be deployed in a containerized "sidecar" configuration with most applications, is configured entirely through environment variables, and all critical functions are exposed via API for configuration via Tofu or post deployment as needed.

## the API

The documentation can be accessed at http://localhost:8888/v1/swagger/index.html.

The logging endpoints are

| HTTP | Endpoint       | Desc                                                         |
| ---- | -------------- | ------------------------------------------------------------ |
| POST | /v1/log        | Log an event with cluster (optional), tags, value (optional) |
| POST | /v1/log/{date} | Log an event with a date (for testing only)                  |

The ETL and analysis endpoints are

| HTTP   | Endpoint                   | Desc                                   |
| ------ | -------------------------- | -------------------------------------- |
| GET    | /v1/etl/{date}/{name}      | Download the contents of an ETL action |
| POST   | /v1/etl/{date}/{name}      | Upload an ETL action                   |
| PUT    | /v1/etl/{date}/{name}      | Run an ETL action                      |
| DELETE | /v1/etl/{date}/{name}      | Remove an ETL action                   |
| POST   | /v1/sequence               | Create an ETL sequence                 |
| GET    | /v1/sequence/{date}/{name} | Run a sequence for a given date        |
| PUT    | /v1/summarize              | Consolidate summary tables             |
| GET    | /v1/summary/{name}         | Fetch the value of the named summary   |

Administrative endpoints include


| HTTP | Endpoint   | Desc                             |
| ---- | ---------- | -------------------------------- |
| GET  | /v1/health | A standard healthcheck endpoint  |
| GET  | /v1/status | Get server stats (RAM, GC, etc.) |

## running its-log

`its-log` is compiled for multiple architectures.

To run on a ARM processors (Mac):

```
make itslog-arm
```

To run on Intel:

```
make itslog-amd
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=navapbc/its-log&type=date&legend=top-left)](https://www.star-history.com/#navapbc/its-log&type=date&legend=top-left)