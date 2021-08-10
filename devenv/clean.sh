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

echo -e "\n\n  run the following to completely remove all cache and temporary build files"
echo -e '    $ rm -rf .cache/ .npm/ ./cmd/frontend/dist /cmd/frontend/node_modules'
