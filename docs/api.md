# The `its-log` API

## logging events

There is one logging endpoint

* POST /v1/log
  
This takes a JSON object of the following shape:

```
type Event struct {
	Cluster   string    `json:"cluster" validate:"max=256"`
	Tags      []string  `json:"tags" validate:"required"`
	Value     string `json:"value" validate:"max=256"`
}
```

Or, by way of example, a minimal event would look like:

```
{
    "tags": [ "v1", "chocolate" ]
}

while an event leveraging all possible features might look like:

```
{
    "cluster": "E56828A0-FEC5-4549-A523-DC762658D6C0",
    "tags": [ "v3", "chocolate", "ramen", "nachos" ],
    "value": "thursday-lunch"
}
```

`its-log` takes responsibility for authentication, as well as attaching a timestamp, the `key_id` of the API key making the request.

## working with the ETL

ETL steps can be added via POST



	// Insert a new ETL step
	auth_adminV1.POST("etl/:date", ETL)
	// Run an ETL step
	auth_adminV1.PUT("etl/:date/:name", ETL)
	// Retrieve the contents of a step, including the last run and run status
	auth_adminV1.GET("etl/:date/:name", ETL)
	// Combine a table from one DB into another DB
	// auth_adminV1.PUT("combine/:source/:destination/:table", Combine)
	auth_adminV1.GET("etl/reload/:date", ReloadEtl)

