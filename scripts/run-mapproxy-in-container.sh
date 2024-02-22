#!/bin/bash

set -e

pushd $(dirname $0)/..

local_port=${MAPPROXY_LISTEN_PORT:-8124}
echo "Will listen on http://localhost:$local_port"

pwd

docker compose build mapproxy
MAPPROXY_LISTEN_PORT=$local_port docker compose up mapproxy
