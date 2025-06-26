#!/bin/bash

set -e

docker build \
    -t hass-update-manager-server \
    -f local/Dockerfile \
    .
