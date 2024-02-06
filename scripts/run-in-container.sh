#!/bin/bash

set -e

image_name="captaincoordinates/tile-id-api"
local_port=${TILE_ID_LISTEN_PORT:-8123}

docker build -t $image_name .
echo "Listening on local port $local_port"
docker run \
    --rm \
    -p $local_port:8080 \
    $image_name
