#!/bin/bash

set -e

pushd $(dirname $0)/..

pip install pre-commit
pre-commit install

curl --fail -o .pre-commit-gofmt.sh https://raw.githubusercontent.com/golang/go/release-branch.go1.1/misc/git/pre-commit
chmod 755 .pre-commit-gofmt.sh

pip install -r scripts/docker/requirements.txt
