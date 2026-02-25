#!/bin/bash

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${scriptDir}/.." || exit 1

export SOPS_AGE_RECIPIENTS=$(<secrets/public-age-keys.txt)
exec 3<<< "$(cat $1)"
sops --encrypt \
    --input-type yaml \
    --output-type yaml \
    --age ${SOPS_AGE_RECIPIENTS} \
    --encrypted-regex "^(encrypted)$" \
    /dev/fd/3
