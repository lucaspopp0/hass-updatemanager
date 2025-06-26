#!/usr/bin/with-contenv bashio

set -x

SWITCHES_JSON=/data/switches.json

if [[ ! -d "/data" ]]; then
    mkdir /data
fi

if [[ ! -f "$SWITCHES_JSON" ]]; then
    echo "{}" > "$SWITCHES_JSON"
fi

/update-manager/server/server
