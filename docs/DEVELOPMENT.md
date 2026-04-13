# development

`its-log` is a [Gin](https://gin-gonic.com/) based Golang API. 

To develop, you should (ideally) configure a debugging environment (provided for VS Code), and for some activities, you will need to run the storage stack, emulating cloud services that `its-log` might rely on in production.

## adding API endpoints

From the "top down," endpoints are attached to the router in `api/endpoints/serve.go`. The function `PourGin` contains a sequence of helpers that attach different endpoints based on functionality.

Using the logging endpoints as an example:

```go
func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *types.Event) {
	auth_logV1 := rG.Group("/")
	permissions := []types.PermissionType{constants.Log, constants.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.POST(constants.LOG_CREATE, LogCreate(ch_evt_out, constants.Log))
	auth_logV1.POST(constants.LOG_CREATE_DATE, LogCreate(ch_evt_out, constants.Test))
}
```

Each endpoint group is granted role-based permissions. The endpoint URLs are defined in `constants`, as they get re-used in testing (and therefore want to be defined in one place). All `gin` handlers have the signature `func (*gin.Context) {}`, but some of the endpoints in `its-log` are [curried](https://wiki.haskell.org/index.php?title=Currying) so that additional parameters can be passed to the handlers (as in the logging portion of the API).

## types

Wherever possible, `its-log` introduces types/abstractions. The internal "its-log time," or `ILTime` structure is an example. Because we move between several time representations (SQLite does not have great time facilities compared to other DB implementations), it is easiest/safest to encapsulate the manipulation of time (and its representation) in one place.

Also, putting base types in one place helps eliminate circular dependencies.

Types are found in `api/internal/types`. 

## e2e

The end-to-end tests are in `api/e2e`. These will stand up the server and execute an entire sequence of operations against the entire API. As of this writing, they require an external S3 service like [ministack](https://ministack.org/), as that is not mocked.

`its-log` prefers E2E tests over unit tests. Where appropriate, adding unit tests is always acceptabe. However, testing happy- and sad-path *business logic* in ways authentic to production use is preferred to testing the correctness of code "out of context."

## etl

`its-log` encapsulates an entire extract/transform/load service. This allows users of `its-log` to write transforms in SQL, Golang, and [Starlark](https://starlark-lang.org/), a Python-like scripting language. These transforms can be added dynamically via API (think as part of "infrastructure as code"), executed, and the results downloaded from the `summary` table.

How to write ETLs in each of these langauges is covered further in the [ETL HOWTO](etl-howto.md).

## adding database queries

We use [sqlc](https://sqlc.dev/) wherever possible for interacting with the database. This is **not** an ORM. Instead, sqlc analyzes SQL code, and then *generates* a Golang wrapper from that SQL. Models for tables are generated, as well as functions for querying and inserting values (as appropriate).

The schema for the `its-log` database is stored in `api/internal/schema/schema.sql`, and the queries we use to search and insert into the datbase are in `api/internal/schema/query.sql`. 

By way of example, here is the query that inserts events into the log table:

```sql
-- name: LogEvent :one
INSERT INTO itslog_events (
  timestamp, key_id, cluster, tags, value
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING id;
```

The annotation is used by sqlc to guide the generation of this code:

```go
type LogEventParams struct {
	Timestamp sql.NullInt64
	KeyID     string
	Cluster   sql.NullString
	Tags      string
	Value     sql.NullString
}

func (q *Queries) LogEvent(ctx context.Context, arg LogEventParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, logEvent,
		arg.Timestamp,
		arg.KeyID,
		arg.Cluster,
		arg.Tags,
		arg.Value,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}
```

Wherever and whenever possible, sqlc should be used for interaction with the database.

> [!NOTE]
> There are currently **no** migrations built into `its-log`. Becasue `its-log` creates a new table each day, we assume (perhaps incorrectly) that we can introduce new models cleanly. However, if versioning is required, adding a key/value `itslog_metadata` table (where we can determine what version of the database is present) may be useful. 
>
> Versioning/migrating the daily databases would likely be a challenge, and may require more substantial architectural thoguht if that day comes to pass.

