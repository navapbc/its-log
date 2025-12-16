# The `its-log` API

There are two endpoints:

* PUT /event/&lt;app-id>/&lt;event><br/>
    Log a timestamped event 

### PUT /event/&lt;app-id>/&lt;event>

Events should be keyed with unique labels.

An event is a string. Internally, `its-log` will map strings to integers (for storage efficiency). Therefore, events should be *consistent*. Good events might look like:

```
PUT //its-log/gov.hhs.cms.bb2/api.patient.v3
PUT //its-log/gov.hhs.cms.bb2/api.patient.v2
PUT //its-log/gov.hhs.cms.bb2/app.start
PUT //its-log/gov.hhs.cms.bfd/api.search_v2
PUT //its-log/gov.hhs.cms.bfd/app.start
```

Each of these events will map internally to a unique value.

Bad events have dynamic elements; this will lead to data that ultimately cannot be analyzed.

For example, here are three event keys that have an time embedded in them.

```
PUT //its-log/gov.hhs.cms.bb2/api.called.2025-12-15-11:34:23
PUT //its-log/gov.hhs.cms.bb2/api.called.2025-12-15-11:34:24
PUT //its-log/gov.hhs.cms.bb2/api.called.2025-12-15-11:34:25
```

These events embed a timestamp, and therefore each would be logged as a separate event. Internally, `its-log` timestamps all data; if the goal is to analyze how many times the API is called, and possibly the time between calls, the correct approach in this case would be to repeatbly log the same event; for example:

```
PUT //its-log/gov.hhs.cms.bb2/api.called
PUT //its-log/gov.hhs.cms.bb2/api.called
PUT //its-log/gov.hhs.cms.bb2/api.called
```
