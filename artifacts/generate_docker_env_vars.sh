#!/usr/bin/env bash

# shellcheck source=artifacts/utils.sh
source "$(dirname "$0")"/utils.sh

verify_var_set 'BRANCH'
verify_var_set 'BETA_BRANCH'

if [ "$BRANCH" = "$BETA_BRANCH" ]
then
  export DOCKER_IMAGE_TAG='beta'
elif [ "$BRANCH" = 'master' ]
then
  export DOCKER_IMAGE_TAG='latest'
else
  echo "[ERROR] Unsupported branch $BRANCH" 1>&2
  exit 1
fi

export DOCKER_IMAGE_NAME="${DOCKER_USERNAME}/${REPOSITORY}:${DOCKER_IMAGE_TAG}"
