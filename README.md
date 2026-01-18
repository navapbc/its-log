# its-log

*It's better than bad, it's good!*

`its-log` is a tiny event logger.

On a Mac M4, `its-log` can sustain logging 30K events/second to a local SQLite database.

`its-log` is an event logger with opinions.

1. Applications generate events
2. Every application must have its own database
3. Every log **could** belong to a *cluster* of logs
4. Every log **must** have a source *category*
5. Every log **should** have an event *category*
6. Every log **could** have a *unique* value

By being opinionated, `its-log` can provide standardized ETL actions that operate over clusters, categories, and values. Users can then develop more specialized ETL actions as needed.

## the API

`its-log` provides interactive Swagger documentation. When run locally with

```
make run
```

The documentation can be accessed at http://localhost:8888/v1/swagger/index.html.

The logging endpoints are

| HTTP   | Endpoint                                    | Desc                                     |
| ------ | ------------------------------------------- | ---------------------------------------- |
| PUT    | /v1/se/{source}/{event}                     | Log a source, event                      |
| PUT    | /v1/sev/{source}/{event}/{value}            | Log a source, event, and value           |
| PUT    | /v1/cse/{cluster}/{source}/{event}          | Log a source and event with a cluster ID |
| PUT    | /v1/csev/{cluster}/{source}/{event}/{value} | ...                                      |

The ETL and analysis endpoints are

| HTTP   | Endpoint                                    | Desc                                     |
| ------ | ------------------------------------------- | ---------------------------------------- |
| GET    | /v1/etl/{date}/{name}                       | Download the contents of an ETL action   |
| POST   | /v1/etl/{date}/{name}                       | Upload an ETL action                     |
| PUT    | /v1/etl/{date}/{name}                       | Run an ETL action                        |
| DELETE | /v1/etl/{date}/{name}                       | Remove an ETL action                     |
| PUT    | /v1/summarize                               | Consolidate summary tables               |
| GET    | /v1/summary/{name}                          | Fetch the value of the named summary     |


## to kick the tires

To run unit tests

```
make test
```

To stand up the API

```
make run
```

To stand up the logger and run the E2E suite

```
make e2e
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=navapbc/its-log&type=date&legend=top-left)](https://www.star-history.com/#navapbc/its-log&type=date&legend=top-left)