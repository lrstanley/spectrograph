#!/bin/bash

if [ "$(basename $PWD)" != "spectrograph" ];then
    >&2 echo "error: run $0 from the root of the project"
    exit 1
fi

mkdir -vp ./cmd/frontend/dist ./cmd/frontend/node_modules/.vite
mkdir -vp ./cmd/http-server/{public,bin}
mkdir -vp ./cmd/worker/bin

DANGLING=$(docker images -f "dangling=true" -q)
if [ ! -z "$DANGLING" ];then
    docker rmi $DANGLING 2>/dev/null
fi

docker-compose \
    --project-name "$(basename $PWD)" \
    --file ./devenv/docker-compose.yml \
    up \
        --always-recreate-deps \
        --abort-on-container-exit \
        --remove-orphans \
        --build \
        --timeout 0
