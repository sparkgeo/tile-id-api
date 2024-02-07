#!/bin/bash

set -e

docker compose build validateopenapi
docker compose run validateopenapi
