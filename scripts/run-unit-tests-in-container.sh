#!/bin/bash

set -e

pushd $(dirname $0)/..

docker compose build api unittest
docker compose run unittest
