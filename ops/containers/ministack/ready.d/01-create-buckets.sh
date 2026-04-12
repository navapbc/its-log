#!/bin/bash

echo "creating default ministack buckets"
apk add --no-cache aws-cli

export ITSLOG_BACKUP_AWS_KEYID="000000000000"
export ITSLOG_BACKUP_AWS_ACCESSKEY="anything"
export AWS_ACCESS_KEY_ID=${ITSLOG_BACKUP_AWS_KEYID}
export AWS_SECRET_ACCESS_KEY=${ITSLOG_BACKUP_AWS_ACCESSKEY}
aws --endpoint-url=http://localhost:4566 s3 mb s3://backup
echo "bucket 000: " $?

export ITSLOG_BACKUP_AWS_KEYID="111111111111"
export ITSLOG_BACKUP_AWS_ACCESSKEY="anything"
export AWS_ACCESS_KEY_ID=${ITSLOG_BACKUP_AWS_KEYID}
export AWS_SECRET_ACCESS_KEY=${ITSLOG_BACKUP_AWS_ACCESSKEY}
aws --endpoint-url=http://0.0.0.0:4566 s3 mb s3://backup
echo "bucket 111: " $?