#!/bin/bash

set -e

docker compose build unittest
docker compose run unittest
