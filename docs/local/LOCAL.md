# running locally

## containerized

### stand up the stack

To start, run the containerized stack:

```
docker compose up --build
```

This runs a stack with 

* its-log
* garage

Garage is a FOSS S3-compatible object server.

### configure garage

Currently, nothing is provided to automatically create buckets
at startup. Therefore, the env must be stood up, logged into, and
configured with at least one bucket. TODO: add automatic config code
at launch time for testing.

```
source configure-garage.env
```

Then, pull the key/secret and put them in `test-secrets.env`

```
source test-secrets.env
```

### reconfigure its-log

Now, its-log needs to know the secrets for the Garage instance. (This could be automated/improved.)

1. Bring down the stack
2. Source the secrets
3. Launch the stack

Now, its-log will have the credentials to talk to the Garage server.

### talk to its-log

This will post messages to its-log.

```
for i in {1..10} ; do http POST localhost:8888/v1/log x-api-key:not-a-real-api-key-but-it-needs-to-be-long source="httpie" event="#hi" value="tacos $i" type="text"; done
```

It should write through to the Garage server.

### check garage

```
aws s3 ls --recursive --endpoint-url http://localhost:3900/ s3://blue-bucket
```

This should list all of the objects in the bucket

```
aws s3 cp --endpoint-url http://localhost:3900/ s3://blue-bucket/2025-12-07/2025-12-07-1765133838-000000001.json .
```

This will yield a file that looks like:

```
{
  "event": "#hi",
  "source": "httpie",
  "type": "text",
  "value": "tacos 1"
}
```