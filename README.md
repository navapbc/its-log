# its-log

*It's better than bad, it's good!*

`its-log` is a tiny event logger.

On a Mac M4, `its-log` can sustain logging 35K events/second to a local SQLite database.

Events are **app-specific, keyed, valueless moments in time**.

* moment in time: an event has a timestamp
* app-specific: each app provides a unique ID to distinguish it from other apps using an its-log instance (e.g. `gov.agency.app`)
* keyed: a key or label tells us what the event was: "v3.api"
* valueless: no value (e.g. `42`) is associated with a key; there are keys but not values.


## the API

There <strike>are two</strike> is one endpoint in the `its-log` API:

* PUT /event/&lt;app-id>/&lt;event><br/>
    Log a timestamped event 

These are described in the [API docs](docs/api.md).

## to kick the tires

```
make serve
```

to run the logger locally, and to run performance tests ("E2E") against the local app:

```
cd k6
k6 run put.js
```
