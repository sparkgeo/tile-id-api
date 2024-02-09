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

cdk deploy --require-approval never
