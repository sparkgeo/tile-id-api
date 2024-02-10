#!/bin/bash

set -e

pushd $(dirname $0)/..

pip install pre-commit
pre-commit install
pip install -r scripts/docker/requirements.txt
