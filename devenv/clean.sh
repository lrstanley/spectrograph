#!/bin/bash

BASE="$(basename $PWD)"

if [ "$BASE" != "spectrograph" ];then
    >&2 echo "error: run $0 from the root of the project"
    exit 1
fi

set -x
docker-compose \
    --project-name "$BASE" \
    --file ./devenv/docker-compose.yml \
    down --volumes --remove-orphans --rmi local --timeout 1

rm -rf .cache/ .npm/
