# its-log

*It's better than bad, it's good!*

`its-log` is a lightweight event logger and ETL environment for resource-constrained, compliance-burdened environments.

On a Mac M4, `its-log` can sustain logging 8K events/second to a local DuckDB database.

## requirements

For local development:

* VSCode
* Golang extensions
* Docker
* `make` and other related tools

## testing

`make e2e`

This will stand up two Docker containers; one that builds and runs `its-log`, and a second that runs an E2E suite against that container.

After running, the test container should report `SUCCESS`, and both containers will exit.

## debugging its-log

A VS Code debugging bundle is prepared and part of this repository, with all of the envrionment variables required to run. Press `F5` to launch `its-log` under the debugger, and `SHIFT-F5` to halt the debugger (on a Mac).

## testing with Postman

A Postman bundle is included under the `docs` directory. It includes a set of scripts to test results, and can serve as a light E2E test. 

## running its-log

`its-log` can be compiled for multiple architectures.

`make build`

This will build the API locally. It is possible to then run the E2E suite locally (after venv creation and requirements installation).

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

| HTTP | Endpoint         | Permissions | Desc                                                                  |
| ---- | ---------------- | ----------- | --------------------------------------------------------------------- |
| POST | /etl/create      | admin       | Create an ETL action                                                  |
| POST | /etl/run         | admin       | Run an ETL action. Takes optional params in the JSON body.            |
| POST | /sequence/create | admin       | Create an ETL sequence                                                |
| POST | /sequence/run    | admin       | Run a sequence for a given date. Takes optional params like /etl/run. |

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

## performance

Performance stress-testing can be run with

`make k6`

This runs K6 in one container, and `its-log` in another.

On an M4 Mac against local NVMe storage, 8K/sec, 4K/sec between containers (HTTP overheads). Cloud performance is likely somewhere in-between.

Containerized performance (2026-10-11)

```
 TOTAL RESULTS 
HTTP
http_req_duration..............: avg=105.4µs min=48.04µs med=98.79µs max=3.99ms p(90)=129.83µs p(95)=147.91µs
  { expected_response:true }...: avg=105.4µs min=48.04µs med=98.79µs max=3.99ms p(90)=129.83µs p(95)=147.91µs
http_req_failed................: 0.00%  0 out of 15600
http_reqs......................: 15600  3970.227365/s

EXECUTION
iteration_duration.............: avg=3.92s   min=3.92s   med=3.92s   max=3.92s  p(90)=3.92s    p(95)=3.92s   
iterations.....................: 1      0.254502/s
vus............................: 1      min=1          max=1
vus_max........................: 1      min=1          max=1

NETWORK
data_received..................: 2.2 MB 548 kB/s
data_sent......................: 3.6 MB 916 kB/s
```



## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=navapbc/its-log&type=date&legend=top-left)](https://www.star-history.com/#navapbc/its-log&type=date&legend=top-left)