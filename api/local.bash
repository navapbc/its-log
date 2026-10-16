#!/usr/bin/env bash

mkdir -p ${HOME}/.its-log/data

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
export ITSLOG_ZAPMODE="debug" # debug or info
export ITSLOG_HASH_SEED=42
export ITSLOG_PROXIES_TRUSTED="TBD"
export ITSLOG_SERVE_HOST="localhost"
export ITSLOG_SERVE_PORT="8888"
export ITSLOG_STORAGE_PATH="${HOME}/.its-log/data"

# export SOPS=$(<../containers/e2e/secrets/itslog-local.yaml)
# echo $SOPS

./its-log serve
