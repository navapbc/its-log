#!/usr/bin/env bash


# These values would be set in a secrets manager
# in a deployment context.
# https://stackoverflow.com/questions/45469133/create-json-file-using-jq
export ITSLOG_APIKEY_PHLLC_ADMIN=$(jq -n \
  --arg app_id "pupper" \
  --arg key_id "pup_admin" \
  --arg permission "admin" \
  --arg key "abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh" \
  '$ARGS.named')
export ITSLOG_APIKEY_PHLLC_LOG=$(jq -n \
  --arg app_id "pupper" \
  --arg key_id "pup_logging" \
  --arg permission "log" \
  --arg key "1234567890123456123456789012345612345678901234561234567890123456" \
  '$ARGS.named')
  
export ITSLOG_BUFFER_FLUSHWAITSEC=1
export ITSLOG_BUFFER_LENGTH=2000
export ITSLOG_GINMODE="debug" # debug or release
export ITSLOG_HASH_SEED=42
export ITSLOG_PROXIES_TRUSTED="TBD"
export ITSLOG_SERVE_HOST="0.0.0.0"
export ITSLOG_SERVE_PORT="8888"
export ITSLOG_STORAGE_PATH="/data"
# A comma-separated list of CIDRs and IP addresses
export ITSLOG_PROXIES="172.18.0.0/12"

export ITSLOG_BACKUP_BUCKET="backup"
export ITSLOG_BACKUP_AWS_ENDPOINT_HOST="ministack"
export ITSLOG_BACKUP_AWS_ENDPOINT_PORT="4566"
export ITSLOG_BACKUP_AWS_ENDPOINT_SCHEME="http"
export ITSLOG_BACKUP_AWS_REGION="us-east-1"
export ITSLOG_BACKUP_AWS_KEYID="000000000000"
export ITSLOG_BACKUP_AWS_ACCESSKEY="anything"

export AWS_ACCESS_KEY_ID=${ITSLOG_BACKUP_AWS_KEYID}
export AWS_SECRET_ACCESS_KEY=${ITSLOG_BACKUP_AWS_ACCESSKEY}

chmod 755 its-log
./its-log serve
