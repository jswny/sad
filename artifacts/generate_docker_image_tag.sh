#!/usr/bin/env bash

if [ "$BRANCH" = "$BETA_BRANCH" ]
then
  DOCKER_IMAGE_TAG='beta'
else
  DOCKER_IMAGE_TAG='latest'
fi

echo "$DOCKER_IMAGE_TAG"
