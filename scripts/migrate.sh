#!/usr/bin/env bash

set -eo pipefail

ROOT="$(dirname "${BASH_SOURCE[0]}")/.."
cd $ROOT
ENV_FILE=".env"

echo "Exporting env vars"
export $(cat $ENV_FILE | grep -v '#' | awk '/=/ {print $1}')

CONTROLLER_HOST=${CONTROLLER_URL%"/v1/graphql"}
cd $ROOT/services/controller

hasura2 migrate apply --all-databases --admin-secret $CONTROLLER_ADMIN_SECRET --endpoint $CONTROLLER_HOST
hasura2 metadata apply --admin-secret $CONTROLLER_ADMIN_SECRET --endpoint $CONTROLLER_HOST