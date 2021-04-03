#!/bin/bash

if [ "$(basename $PWD)" != "spectrograph" ];then
    >&2 echo "error: run $0 from the root of the project"
    exit 1
fi

docker-compose \
    --project-name "$(basename $PWD)" \
    --file ./devenv/docker-compose.yml \
    rm --stop -v --force
