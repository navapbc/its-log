# its-log

*It's better than bad, it's good!*

`its-log` is a tiny event logger.

On a Mac M4, `its-log` can sustain logging 30K events/second to a local SQLite database.

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

To stand up the API

```
make up
```

To stand up the logger and run the E2E suite

```
make e2e
```

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=navapbc/its-log&type=date&legend=top-left)](https://www.star-history.com/#navapbc/its-log&type=date&legend=top-left)