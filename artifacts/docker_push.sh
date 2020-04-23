#!/usr/bin/env bash

# shellcheck source=artifacts/utils.sh
source "$(dirname "$0")"/utils.sh

verify_var_set 'DOCKER_PASSWORD'
verify_var_set 'DOCKER_USERNAME'

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

verify_var_set 'REPOSITORY'

docker build --tag "$REPOSITORY" .

verify_var_set 'DOCKER_IMAGE_NAME'

docker tag "$REPOSITORY" "$DOCKER_IMAGE_NAME"

docker push "$DOCKER_IMAGE_NAME"
