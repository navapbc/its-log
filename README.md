# its-log

*It's better than bad, it's good!*

`its-log` is a tiny event logger.

Events are **app-specific, keyed, valueless moments in time**.

* moment in time: an event has a timestamp
* app-specific: each app provides a unique ID to distinguish it from other apps using an its-log instance (e.g. `gov.hhs.cms.bb2`)
* keyed: a key or label tells us what the event was: "v3.patient.api"
* valueless: no value (e.g. `42`) is associated with a key; there are keys but not values.

## the API

There are two endpoints in the `its-log` API:

* PUT /event/&lt;app-id>/&lt;event><br/>
    Log a timestamped event 
* PUT /unique/&lt;app-id>/&lt;event><br/>
    Log an event today, uniquely

These are described in the [API docs](docs/api.md).

## packages used

* [github.com/spf13/viper](https://github.com/spf13/viper)
* go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
* go get gocloud.dev

aws s3 ls s3://blue-bucket/ --endpoint-url http://localhost:3900
