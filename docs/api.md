# The `its-log` API

There are two endpoints:

* PUT /event/&lt;app-id>/&lt;event><br/>
    Log a timestamped event 
* PUT /unique/&lt;app-id>/&lt;event><br/>
    Log an event today, uniquely

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

### PUT /unique/&lt;app-id>/&lt;event>

Unique events are only logged once per day. Using *limited* dynamic event tags makes sense for unique events.

For example, 

```
PUT //its-log/gov.hhs.cms.bb2/auth.app.alicemedical
PUT //its-log/gov.hhs.cms.bb2/auth.app.bobshealth
PUT //its-log/gov.hhs.cms.bb2/auth.app.bobshealth
```

This sequence will yield two events in the database: one for the event `auth.app.alicemedical` and one for `auth.app.bobshealth`. Only the first event for each tag will be stored. Logs rotate daily; therefore, if the same sequence comes in on the next day, again, only two events will be logged.

Note that CUI/PHI/PII should not be used for events. If (say) app names are considered "sensitive," then the *application* should do something like the following:

1. Add a secret to your environment; call this `LOGGING_SALT`. A `sha1(uuid.uuid4()).hexdigest()` is suitable
2. `app_name_salted = f'{app_name}:{LOGGING_SALT}'`
3. Log a key using the salted name; for example, `auth.app.{app_name_salted}`

Now, events should take the form:

```
PUT its-log/gov.hhs.cms.bb2/auth.app.031edd7d41651593c5fe5...
```

This will allow for analysis that (for example) can assert that there were 36 unique applications logged in during a given day, or if the same application logged in every day for a week. Using this approach, it is *not* possible to say (from the data stored in `its-log`) that the application logging in every day was `bobshealth`.
