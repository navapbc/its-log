#!/bin/bash

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${scriptDir}/../../secrets" || exit 1

export SOPS_AGE_RECIPIENTS=$(<public-age-keys.txt)
exec 3<<< "$(cat $1.sops.yaml)"
sops --encrypt \
    --input-type yaml \
    --output-type yaml \
    --age ${SOPS_AGE_RECIPIENTS} \
    --encrypted-regex "/sensitive/" \
    /dev/fd/3 > $1.sops.enc.yaml
popd