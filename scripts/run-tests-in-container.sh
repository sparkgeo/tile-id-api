#!/bin/bash

set -e

pushd $(dirname $0)/..

scripts/run-validate-openapi-in-container.sh
scripts/run-unit-tests-in-container.sh
scripts/run-integration-tests-in-container.sh
