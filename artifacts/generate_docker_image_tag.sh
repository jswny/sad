#!/usr/bin/env bash

if [ "$BRANCH" = "$BETA_BRANCH" ]
then
  DOCKER_IMAGE_TAG='beta'
elif [ "$BRANCH" = 'master' ]
then
  DOCKER_IMAGE_TAG='latest'
else
  echo "ERROR: unsupported branch $BRANCH" 1>&2
  exit 1
fi

echo "$DOCKER_IMAGE_TAG"
