#!/bin/bash

if [ "$(basename $PWD)" != "spectrograph" ];then
    >&2 echo "error: run $0 from the root of the project"
    exit 1
fi

mkdir -p .cache/{frontend_dist,node_dotnpm,node_modules,go_cache,go_path}
mkdir -p .npm/
# mkdir -vp ./cmd/frontend/dist ./cmd/frontend/node_modules/.vite
mkdir -p ./cmd/http-server/{public,bin}
mkdir -p ./cmd/worker/bin

DANGLING=$(docker images -f "dangling=true" -q)
if [ ! -z "$DANGLING" ];then
    docker rmi $DANGLING 2>/dev/null
fi

export USER=$(id -u)
export GROUP=$(id -g)

docker-compose \
    --project-name "$(basename $PWD)" \
    --file ./devenv/docker-compose.yml \
    up \
        --always-recreate-deps \
        --abort-on-container-exit \
        --remove-orphans \
        --build \
        --timeout 0
