#!/bin/bash

set -e

docker compose build api unittest
docker compose run unittest
