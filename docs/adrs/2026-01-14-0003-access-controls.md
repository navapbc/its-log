# Access controls and API keys

Date: 2025-01-14

## Status

Accepted

## Context

`its-log` is designed to support multiple applications logging events at the same time. However, we want to keep the data from different applications separate. Further, we want to keep different *kinds* of operations separate: logging is different from loading ETL steps or invoking the ETL process, for example.

Therefore, `its-log` will bind critical metadata (e.g. application ids) to API keys, so that the shared secret (the key) can be used to determine what application, with what level of permission, is accessing the `its-log` API.

## Decision

`its-log` manages access via symmetric secrets. In a secrets manager, we will store an API key object that looks like:

```
{
    "app_id": "alices_app", 
    "key_id": "alice_log_key", 
    "permission": "log", 
    "key": 
    "abcdefghabcdefghabcdefghabcdefgh"
}
```

This secret tells us what application this key is for (the `app_id`), the unique id for this key (`key_id`), the permission(s) granted this key (`log`, `admin`, `test`, etc.), and the symmetric secret itself. The key must be at least 32 bytes long and achieve minimal entropy levels. 

```
openssl rand -base64 48
```

or similar, with an appropriate crypto implementation, is sufficient. These secrets are then provided to the `its-log` server at launch time as environment variables.

Under local development circumstances, the secrets would look like:

```
"ITSLOG_APIKEY_ALICE": "{\"app_id\": \"alices_app\", \"key_id\": \"alice_log_key\", \"permission\": \"log\", \"key\": \"abcdefghabcdefghabcdefghabcdefgh\"}",
"ITSLOG_APIKEY_BOB": "{\"app_id\": \"bobs_app\", \"key_id\": \"bob_log_key\", \"permission\": \"log\", \"key\": \"12345678901234561234567890123456\"}",
```

These values need to be shared back with the owner of the application who will be using the `its-log` logger.

`its-log` creates a separate database for each application; this is what the `app_id` is for. When an API call is made, `its-log` uses the API key to look up the `app_id`; in this way, it is not possible for one app to log to another application's database. Further, each of the endpoint groups have different permissions associated with them. For example, a key that has the permission `log` is not able to load or execute the ETL pipeline. 

## Consequences

A secret must be shared somewhere for an application to use `its-log`. We choose to use a symmetric API key and associate metadata with that key. A single user of `its-log` may end up with multiple keys, one for logging, one for CI/CD, and one for testing in a lower environment (for example).

At this time, there is no administrative interface to dynamically add/remove keys and permissions. This operation must be performend manually within the secrets manager used by `its-log`. It is conceivable that this could be worked into infrastructure as code (e.g. with local secrets encryption), but that is a deployment detail that is beyond the scope of this ADR.