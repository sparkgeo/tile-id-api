#!/bin/bash

set -e

pushd $(dirname $0)/..

docker compose build validateopenapi
docker compose run validateopenapi
