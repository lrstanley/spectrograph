#!/bin/sh

function kill_descendant_processes {
    local pid="$1"
    local and_self="${2:-false}"
    if children="$(pgrep -P "$pid")"; then
        for child in $children; do
            kill_descendant_processes "$child" true
        done
    fi
    if [[ "$and_self" == true ]]; then
        kill -TERM "$pid"
    fi
}

make debug &
make generate-watch &

wait -n
kill_descendant_processes $$
wait
echo -e "\n\n"