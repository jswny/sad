#!/bin/bash

# Exit on any error, undefined variable, or pipe failure 
set -euo pipefail

verify_var_set() {
  if [ -z "${!1}" ]; then
    if [ -z "$2" ]; then
      log 'error' "\"$1\" is blank or unset!"
    else
      log 'error' "$2"
    fi
    exit 1
  fi
}

# Translate input environment variables
deploy_server="$INPUT_DEPLOY_SERVER"
deploy_username="$INPUT_DEPLOY_USERNAME"
deploy_password="$INPUT_DEPLOY_PASSWORD"
deploy_root_dir="$INPUT_DEPLOY_ROOT_DIR"
encrypted_deploy_key_cypher_key="$INPUT_ENCRYPTED_DEPLOY_KEY_CYPHER_KEY"
encrypted_deploy_key_iv="$INPUT_ENCRYPTED_DEPLOY_KEY_IV"
app_path="$INPUT_PATH"
debug="$INPUT_DEBUG"

ssh_key_types='rsa,dsa,ecdsa'

verify_var_set 'deploy_server'
verify_var_set 'deploy_username'
verify_var_set 'deploy_password'
verify_var_set 'deploy_root_dir'
verify_var_set 'encrypted_deploy_key_cypher_key'
verify_var_set 'encrypted_deploy_key_cypher_iv'
verify_var_set 'app_path' 'path is blank or unset!'
verify_var_set 'debug'
verify_var_set 'ssh_key_types'

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

cat "/github/workspace/${app_path}/.gitignore"

local_image_id="$(docker images -q "$GITHUB_REPOSITORY" 2> /dev/null)"

verify_var_set 'local_image_id' 'No local Docker image detected for this repository! Please build a local image first before deploying!'

log 'debug' "Local Docker image ID: $local_image_id"

local_image="$(docker inspect --format='{{ (index .RepoTags 0) }}' $local_image_id 2> /dev/null)"

verify_var_set 'local_image' 'Could not find the local Docker image name and tag!'

log 'debug' "Local Docker image name and tag: $local_image"

openssl aes-256-cbc -K "$encrypted_deploy_key_cypher_key" -iv "$encrypted_deploy_key_iv" -in "$app_path/deploy_key.enc" -out "deploy_key" -d

chmod 600 'deploy_key'

ssh-add 'deploy_key'

{ ssh-keyscan -t "$ssh_key_types" -H "$deploy_server" >> "${HOME}/.ssh/known_hosts"; } 2>&1
