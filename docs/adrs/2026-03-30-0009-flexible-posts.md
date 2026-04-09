# Flexible POST requests for ETLs

Date: 2026-03-30

## Status

Accepted

## Context

There are currently two classes of ETL action that can be taken in `its-log`.

1. SQL
2. Golang

### SQL ETL actions

The SQL ETL actions are either defaulted in, or can be loaded via API. Either way, they're only able to execute within the context of the single SQLite database they are part of. This means they are most appropriate for analyses of (say) the events table, and the generation of summary information in the summary table. The only external context that can be passed into an SQL ETL action (at this time) is the `key_id`, the `app_id`, and the `date`. 

For example: an SQL ETL action might be to count the total number of events in the `events` table, and insert that count in the `summary` table with an opertaion called `count.total`. 

### Golang ETL actions

Golang ETL actions are *different*. These ETL actions are *not* dynamic; that is, Golang ETL actions **must** be compiled into `its-log`. 

These ETLs are Go functions that match the following signature:

```
func (etlP *types.RunEtlParams) error
```

where `RunEtlParams` are:

```
type RunEtlParams struct {
	AppId   string
	GinCtx  *gin.Context
	EtlName string
	KeyId   string
	Storage *Storage
	Payload map[string]any
}
```

This means a Golang ETL action can:

* write to the Gin context (e.g. write a response to the client---but shouldn't)
* the name of the ETL being executed
* the AppId and KeyId as looked up based on the API key used
* a link to the current (open) logging database file
* a payload of arbitrary JSON data

The last element is part of what makes them flexible; when calling a Golang ETL directly (or as part of a sequence), JSON data can be passed as part of the POST to feed the ETL. For example, the `Consolidate` ETL action expects to be provided with a list of dates. These dates let us open older/prior databases and copy their summary data forward to the current database. That list is passed as the JSON parameter `prior_summary_dates_to_include`. (The details may change.) In this way, when calling the `Consolidate` action, the user of the API can specify that they want to copy forward:

1. Just yesterday
2. The past week
3. The past month
4. The past year

The intent is that on 2026-03-30, the data from 2026-03-29 will be copied forward... and on the 29th, the data from the 28th will be copied forward, and on the 28th, the data from the... the base case is that there is no prior day, but ultimately, the most recent `summary` table will always have contained within it all of the summaries since `its-log` began operation. (This particular model may change as well. The point of this example is how a Golang ETL action can work with flexible POST parameters, not the commitment to this particular data copying model.)

It is in this way that Golang ETLs can be used to export CSV files, copy data to other systems (e.g. S3, Postgres), and so forth. Because they are arbitrary Golang functions, they can prepare, move, and manipulate data in ways that pure SQL cannot. And, because we can pass (non-conflicting) JSON keys to the POST calls that drive both individual ETL and sequenced ETL actions, they can be parameterized to meet the needs of multiple systems.

## Decision

Golang ETLs will leverage flexibly-defined POST parameters on a per-ETL-function basis.
