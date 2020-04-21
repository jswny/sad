#!/usr/bin/env bash

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker build --tag "$REPOSITORY" .
docker tag "$REPOSITORY" "$DOCKER_IMAGE_NAME"
docker push "$DOCKER_IMAGE_NAME"