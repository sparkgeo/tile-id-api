#!/bin/bash

set -e

local_port=${TILE_ID_LISTEN_PORT:-8123}
echo "Will listen on http://localhost:$local_port"

docker compose build api
TILE_ID_LISTEN_PORT=$local_port docker compose up api
