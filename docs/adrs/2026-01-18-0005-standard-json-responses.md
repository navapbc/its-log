# Standard API responses

Date: 2026-01-18

## Status

Accepted

## Context

JSON responses from `its-log` should follow a standard pattern.

## Decision

Prometheus is a widely used, open source metrics application.

`its-log` will use the Prometheus JSON response format when responding to all API calls.

From [the Prometheus documentation](https://prometheus.io/docs/prometheus/latest/querying/api/):

```
{
  "status": "success" | "error",
  "data": <data>,

  // Only set if status is "error". The data field may still hold
  // additional data.
  "errorType": "<string>",
  "error": "<string>",

  // Only set if there were warnings while executing the request.
  // There will still be data in the data field.
  "warnings": ["<string>"],
  // Only set if there were info-level annotations while executing the request.
  "infos": ["<string>"]
}
```

## Consequences

On successful calls, `its-log` will return the minimal response of

```
{ "status": "success" }
```

keeping the number of bytes shipped low. In the case of warnings and errors, a richer JSON response will be provided as described above.

In following this format, documenting the API will be simpler (because results are consisten), and we follow the practices of a widely-used, open source system that is similar in kind and character.