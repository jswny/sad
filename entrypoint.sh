#!/bin/sh

# Exit on any error, undefined variable, or pipe failure 
set -euo pipefail

# Translate input environment variables
deploy_server="$INPUT_DEPLOY_SERVER"
deploy_username="$INPUT_DEPLOY_USERNAME"
deploy_password="$INPUT_DEPLOY_PASSWORD"
deploy_root_dir="$INPUT_DEPLOY_ROOT_DIR"
encrypted_deploy_key_cypher_key="$INPUT_ENCRYPTED_DEPLOY_KEY_CYPHER_KEY"
encrypted_deploy_key_cypher_iv="$INPUT_ENCRYPTED_DEPLOY_KEY_CYPHER_IV"
app_path="$INPUT_PATH"
debug="$INPUT_DEBUG"

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

verify_var_set() {
  if [ -z "${!1}" ]; then
    log 'error' "$1 is blank or unset!"
    exit 1
  fi
}

cat "/github/workspace/$app_path/.gitignore"

docker images

LOCAL_IMAGE_ID="$(docker images -q "$GITHUB_REPOSITORY" 2> /dev/null)"

if [[ "$LOCAL_IMAGE_ID" == "" ]]; then
  log 'error' 'No local Docker image detected for this repository! Please build a local image first before deploying!'
  exit 1
fi

log 'debug' "Local Docker image ID: $LOCAL_IMAGE_ID"
