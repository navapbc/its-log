#!/usr/bin/env bash

# These values would be set in a secrets manager
# in a deployment context.

export ITSLOG_APIKEY_PHLLC_ADMIN="{\"app_id\": \"pupper_health_llc\", \"key_id\": \"ph_admin\", \"permission\": \"admin\", \"key\": \"abcdefghabcdefghabcdefghabcdefgh\"}"
export ITSLOG_APIKEY_PHLLC_LOG="{\"app_id\": \"pupper_health_llc\", \"key_id\": \"ph_logging\", \"permission\": \"log\", \"key\": \"12345678901234561234567890123456\"}"

export ITSLOG_BUFFER_FLUSHWAITSEC=1
export ITSLOG_BUFFER_LENGTH=2000
export ITSLOG_GINMODE="debug" # debug or release
export ITSLOG_HASH_SEED=42
export ITSLOG_PROXIES_TRUSTED="TBD"
export ITSLOG_SERVE_HOST="localhost"
export ITSLOG_SERVE_PORT="8888"
export ITSLOG_STORAGE_PATH="${PWD}/data/itslog"

./api/its-log serve
