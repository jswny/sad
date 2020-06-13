#!/bin/sh

# Exit on any error, undefined variable, or pipe failure 
set -euo pipefail

log() {
  local prefix_spacer="-----"
  local component="Deploy"
  local prefix="$prefix_spacer [$component]"
  if [ "$1" = 'debug' ]; then
    if [ "$debug" = 1 ]; then
      echo "$prefix DEBUG: $2"
    fi
  elif [ "$1" = 'info' ]; then
    echo "$prefix INFO: $2"
  elif [ "$1" = 'warn' ]; then
    echo "$prefix WARN: $2" >&2
  elif [ "$1" = 'error' ]; then
    echo "$prefix ERROR: $2" >&2
  else
    echo "$prefix INTERNAL ERROR: invalid option \"$1\" for log() with message \"$2\"" >&2 "" >&2
  fi
}

# Translate input environment variables
DEPLOY_SERVER="$INPUT_DEPLOY_SERVER"
DEPLOY_USERNAME="$INPUT_DEPLOY_USERNAME"
DEPLOY_PASSWORD="$INPUT_DEPLOY_PASSWORD"
DEPLOY_ROOT_DIR="$INPUT_DEPLOY_ROOT_DIR"
ENCRYPTED_DEPLOY_KEY_CYPHER_KEY="$INPUT_ENCRYPTED_DEPLOY_KEY_CYPHER_KEY"
ENCRYPTED_DEPLOY_KEY_CYPHER_IV="$INPUT_ENCRYPTED_DEPLOY_KEY_CYPHER_IV"
APP_PATH="$INPUT_PATH"
DEBUG="$INPUT_DEBUG"

cat "/github/workspace/$APP_PATH/.gitignore"

docker images

LOCAL_IMAGE_ID="$(docker images -q "$GITHUB_REPOSITORY" 2> /dev/null)"

if [[ "$LOCAL_IMAGE_ID" == "" ]]; then
  log 'error' 'No local Docker image detected for this repository! Please build a local image first before deploying!'
  exit 1
fi

log 'debug' "Local Docker image ID: $LOCAL_IMAGE_ID"


