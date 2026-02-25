#!/bin/bash

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
pushd "${scriptDir}/../../secrets" || exit 1

export SOPS_AGE_KEY_FILE=age-key.txt
exec 3<<< "$(cat $1.sops.enc.yaml)"
sops --decrypt \
    --input-type yaml \
    --output-type yaml \
    /dev/fd/3
popd