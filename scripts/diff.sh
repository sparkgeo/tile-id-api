#!/bin/bash

set -e

pushd $(dirname $0)/../iac

if [ -z "$AWS_REGION" ]; then
    echo "AWS_REGION env var is required for deployment"
    exit 1
fi
if [ -z "$AWS_ACCOUNT" ]; then
    echo "AWS_ACCOUNT env var is required for deployment"
    exit 1
fi

cdk diff \
    -c aws_region=$AWS_REGION \
    -c aws_account=$AWS_ACCOUNT \
    --require-approval never
