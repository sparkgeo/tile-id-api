#!/bin/bash

# =====
# Should only be deployed from local environment from main branch.
# Main branch GitHub Actions must have completed successfully first.
# Best practice would be to deploy from CICD, e.g. via GitHub Actions,
# but as this is a personal repo it does not have access to Sparkgeo's
# organizational secrets for AWS access.
# Login locally via AWS SSO before executing this script.
# =====

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

cdk deploy \
    -c aws_region=$AWS_REGION \
    -c aws_account=$AWS_ACCOUNT \
    --require-approval never
